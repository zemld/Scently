from pathlib import Path
import requests
from bs4 import BeautifulSoup
from concurrent.futures import ThreadPoolExecutor, as_completed
from url_collectors.check_url import is_perfume_url, is_product_url

OUTPUT_DIR = Path.cwd() / "collected_urls"


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


def get_urls_from_sitemap(sitemap: str, limit: int | None = None) -> list[str]:
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
    # TODO: при развертывании приложения нужно не забыть убрать ограничение
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
        futs = {ex.submit(is_perfume_url, u): u for u in urls.to_process}
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
