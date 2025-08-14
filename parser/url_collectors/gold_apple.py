import time
import re
from pathlib import Path
import requests
from requests import HTTPError
from bs4 import BeautifulSoup
from concurrent.futures import ThreadPoolExecutor, as_completed

PRODUCT_URL_RE = re.compile(r"^https://goldapple\.ru/\d{6,}-[a-z0-9-]+$", re.IGNORECASE)
OUTPUT_DIR = Path.cwd() / "collected_urls"

_PERFUME = [
    re.compile(r"\bпарфюмированная\s+вода\b", re.IGNORECASE),
    re.compile(r"\bпарфюмерная\s+вода\b", re.IGNORECASE),
    re.compile(r"\bтуалетная\s+вода\b", re.IGNORECASE),
    re.compile(r"\bэкстракт\s+духов\b", re.IGNORECASE),
    re.compile(r"\bдухи\b", re.IGNORECASE),
    re.compile(r"\beau\s*de\s*parfum\b", re.IGNORECASE),
    re.compile(r"\beau\s*de\s*toilette\b", re.IGNORECASE),
    re.compile(r"\beau\s*de\s*cologne\b", re.IGNORECASE),
    re.compile(r"\bEDP\b", re.IGNORECASE),
    re.compile(r"\bEDT\b", re.IGNORECASE),
    re.compile(r"\bEDC\b", re.IGNORECASE),
]


def _get_output_file(index: int) -> Path:
    return OUTPUT_DIR / f"gold-apple-{index}.txt"


def normalize_href(href: str) -> str | None:
    if not href:
        return None
    if href.startswith("//"):
        href = "https:" + href
    elif href.startswith("/"):
        href = "https://goldapple.ru" + href
    href = href.split("?", 1)[0].rstrip("/")
    return href


def is_product_url(href: str) -> bool:
    return bool(PRODUCT_URL_RE.match(href))


def _is_perfume_html(html: str) -> bool:
    soup = BeautifulSoup(html, "lxml")
    return any(rx.search(soup.title.string.strip()) for rx in _PERFUME)


class PerfumeResponse:
    url: str
    is_perfume: bool
    is_retry_needed: bool


def _is_perfume_url(url: str, timeout: int = 12) -> PerfumeResponse:
    time.sleep(1)
    answer = PerfumeResponse()
    answer.url = url
    answer.is_retry_needed = False
    try:
        r = requests.get(url, timeout=timeout, headers={"User-Agent": "Mozilla/5.0"})
        r.raise_for_status()
        answer.is_perfume = _is_perfume_html(r.content)
        return answer
    except HTTPError as e:
        print(f"HTTPError: {e}")
        answer.is_retry_needed = e.response.status_code == 429
        answer.is_perfume = False
    except Exception:
        return False


def get_sitemaps() -> list[str]:
    headers = {"User-Agent": "Mozilla/5.0"}

    try:
        r = requests.get(
            "https://goldapple.ru/sitemap.xml", timeout=20, headers=headers
        )
        r.raise_for_status()
        soup = BeautifulSoup(r.content, "xml")

        sitemaps = [sitemap.find("loc").string for sitemap in soup.find_all("sitemap")]
        sitemaps = [sitemap for sitemap in sitemaps if sitemap]
        return sitemaps
    except Exception:
        return []


def get_urls_from_sitemap(sitemap: str, limit: int) -> list[str]:
    headers = {"User-Agent": "Mozilla/5.0"}

    try:
        r = requests.get(sitemap, timeout=30, headers=headers)
        r.raise_for_status()
        soup = BeautifulSoup(r.content, "xml")
        urls = [url.string.strip() for url in soup.find_all("loc")]
        return urls[:limit] if limit else urls
    except Exception as e:
        return []


def get_product_urls_from_sitemap(sitemap: str) -> list[str]:
    urls = get_urls_from_sitemap(sitemap, 1000)
    urls = [normalize_href(u) for u in urls]
    urls = [u for u in urls if u and is_product_url(u)]
    return urls


class ProcessedUrls:
    perfume_urls: set[str]
    to_process: list[str]

    def __init__(self, perfume_urls: set[str] = set(), to_process: list[str] = []):
        self.perfume_urls = perfume_urls
        self.to_process = to_process


def process_urls(urls: ProcessedUrls, max_workers: int = 16) -> ProcessedUrls:
    processed = ProcessedUrls()
    print(f"Processing {len(urls.to_process)} URLs...")
    with ThreadPoolExecutor(max_workers=max_workers) as ex:
        futs = {ex.submit(_is_perfume_url, u): u for u in urls.to_process}
        for fut in as_completed(futs):
            u = futs[fut]
            try:
                result = fut.result()
                if result.is_perfume:
                    processed.perfume_urls.add(u)
                elif result.is_retry_needed:
                    processed.to_process.append(u)
            except Exception:
                pass
    print(f"Processed {len(processed.perfume_urls)} perfume URLs.")
    return processed


def get_perfume_urls_from_sitemap(sitemap: str, max_workers: int = 16) -> set[str]:
    product_urls = get_product_urls_from_sitemap(sitemap)
    if not product_urls:
        return set()

    processed = ProcessedUrls(to_process=product_urls)

    while processed.to_process:
        processed = process_urls(processed, max_workers=max_workers)

    return processed.perfume_urls


def collect_links_from_sitemap(max_workers: int = 16) -> None:
    sitemaps = get_sitemaps()
    for i, sitemap in enumerate(sitemaps):
        urls = get_perfume_urls_from_sitemap(sitemap)
        if urls:
            print(f"Collected {len(urls)} URLs from {sitemap}.")
            _get_output_file(i).write_text(
                "\n".join(sorted(urls)) + "\n", encoding="utf-8"
            )


if __name__ == "__main__":
    collect_links_from_sitemap()
