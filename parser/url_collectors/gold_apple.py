import re
from pathlib import Path
import requests
from xml.etree import ElementTree as ET
from concurrent.futures import ThreadPoolExecutor, as_completed

PRODUCT_URL_RE = re.compile(r"^https://goldapple\.ru/\d{6,}-[a-z0-9-]+$", re.IGNORECASE)
SM_NS = {"sm": "http://www.sitemaps.org/schemas/sitemap/0.9"}
OUTPUT_FILE = Path.cwd() / "collected_urls" / "gold_apple.txt"

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
    text = html.lower()
    m = re.search(r"<title>(.*?)</title>", text, re.I | re.S)
    hay = m.group(1) if m else text
    pos = any(rx.search(hay) for rx in _PERFUME)
    if pos:
        return True
    pos2 = any(rx.search(text) for rx in _PERFUME)
    return pos2


def _is_perfume_url(url: str, timeout: int = 12) -> bool:
    try:
        r = requests.get(url, timeout=timeout, headers={"User-Agent": "Mozilla/5.0"})
        if not r.ok:
            return False
        return _is_perfume_html(r.text)
    except Exception:
        return False


def get_sitemaps() -> list[str]:
    headers = {"User-Agent": "Mozilla/5.0"}

    try:
        r = requests.get(
            "https://goldapple.ru/sitemap.xml", timeout=20, headers=headers
        )
        r.raise_for_status()
        root = ET.fromstring(r.content)

        sitemaps = [
            e.find("sm:loc", SM_NS).text
            for e in root.findall("sm:sitemap", SM_NS)
            if e.find("sm:loc", SM_NS) is not None
        ]

        return sitemaps
    except Exception:
        return []


def get_urls_from_sitemap(sitemap: str) -> list[str]:
    headers = {"User-Agent": "Mozilla/5.0"}

    try:
        r = requests.get(sitemap, timeout=30, headers=headers)
        r.raise_for_status()
        root = ET.fromstring(r.content)

        urls = [
            u.find("sm:loc", SM_NS).text
            for u in root.findall("sm:url", SM_NS)
            if u.find("sm:loc", SM_NS) is not None
        ]

        return urls
    except Exception:
        return []


def get_product_urls_from_sitemap(sitemap: str) -> list[str]:
    urls = get_urls_from_sitemap(sitemap)
    urls = [normalize_href(u) for u in urls]
    urls = [u for u in urls if u and is_product_url(u)]
    return urls


def get_perfume_urls_from_sitemap(sitemap: str, max_workers: int = 16) -> set[str]:
    product_urls = get_product_urls_from_sitemap(sitemap)
    if not product_urls:
        return set()

    product_links: set[str] = set()
    with ThreadPoolExecutor(max_workers=max_workers) as ex:
        futs = {ex.submit(_is_perfume_url, u): u for u in product_urls}
        for fut in as_completed(futs):
            u = futs[fut]
            try:
                if fut.result():
                    product_links.add(u)
            except Exception:
                pass

    return product_links


def collect_links_from_sitemap(max_workers: int = 16) -> set[str]:
    product_links: set[str] = set()
    sitemaps = get_sitemaps()
    for sitemap in sitemaps:
        perfume_set = get_perfume_urls_from_sitemap(sitemap, max_workers=max_workers)
        product_links.update(perfume_set)

    return product_links


if __name__ == "__main__":
    links = collect_links_from_sitemap(max_workers=24)
    if links:
        OUTPUT_FILE.write_text("\n".join(links), encoding="utf-8")
