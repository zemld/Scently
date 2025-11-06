import logging
import re
import time
from collections.abc import Iterable, Sequence
from concurrent.futures import ThreadPoolExecutor, as_completed
from threading import Lock
from urllib.parse import parse_qs, urlencode, urljoin, urlparse

from bs4 import BeautifulSoup
from tqdm import tqdm

from models import Perfume
from util import get_page

from ..page_parser import PageParser
from ..scrapper import Scrapper

_CATALOG_BASE_URL = "https://goldapple.ru/parfjumerija"
_DEFAULT_PRODUCT_TYPE_IDS = [245, 639, 640, 989, 1108]
_PRODUCT_URL_RE = re.compile(
    r"^https://goldapple\.ru/(?:.*?\D)?\d{5,}(?:\D.*)?$",
    re.IGNORECASE,
)


def _collect_product_links(
    product_type_ids: Sequence[int] | Iterable[int],
    *,
    max_pages: int | None = None,
) -> list[str]:
    product_ids_str = _normalize_product_type_ids(product_type_ids)
    if not product_ids_str:
        raise ValueError("product_type_ids must contain at least one id")

    collected_links: list[str] = []
    seen_links: set[str] = set()

    current_page = 1
    last_page = max_pages

    while True:
        print(f"Collecting product links from page {current_page}...")
        if max_pages is not None and current_page > max_pages:
            break

        page_url = _build_catalog_url(product_ids_str, current_page)

        retries = 0
        backoff_seconds = 2.0
        page = None
        while retries < 4:
            page = get_page(page_url, use_playwright=True)
            if page:
                break
            retries += 1
            time.sleep(backoff_seconds)
            backoff_seconds *= 2
        if not page:
            print(
                f"Failed to load catalog page {current_page} after retries; skipping.",
            )
            current_page += 1
            continue

        page_links = _extract_product_links(page)
        if not page_links:
            time.sleep(1.0)
            retry_page = get_page(page_url, use_playwright=True)
            if retry_page:
                page_links = _extract_product_links(retry_page)
        if not page_links:
            logging.info(
                "No product links found on page %s; assuming end of catalog.",
                current_page,
            )
            break
        print(f"Found {len(page_links)} product links on page {current_page}.")
        new_links = 0
        for link in page_links:
            if link in seen_links:
                continue
            seen_links.add(link)
            collected_links.append(link)
            new_links += 1

        if new_links == 0:
            break

        detected_last_page = _detect_last_page(page)
        if detected_last_page is not None:
            if last_page is None:
                last_page = detected_last_page
            else:
                last_page = max(last_page, detected_last_page)

        current_page += 1

        if last_page is not None and current_page > last_page:
            break

    return collected_links


def _normalize_product_type_ids(
    product_type_ids: Sequence[int] | Iterable[int],
) -> str:
    normalized_ids: list[str] = []
    for product_type_id in product_type_ids:
        text_id = str(product_type_id)
        if text_id:
            normalized_ids.append(text_id)
    return ",".join(normalized_ids)


def _build_catalog_url(product_type_ids: str, page: int) -> str:
    params: dict[str, str | int] = {"producttype": product_type_ids}
    if page > 1:
        params["p"] = page
    return f"{_CATALOG_BASE_URL}?{urlencode(params)}"


def _extract_product_links(page: BeautifulSoup) -> list[str]:
    product_links: list[str] = []
    for card in page.select("article[data-scroll-id]"):
        anchor = card.find("a", href=True)
        if not anchor:
            for attr in ("data-url", "data-href", "data-product-url", "data-link"):
                data_val = card.get(attr)
                if not isinstance(data_val, str):
                    continue
                normalized = urljoin(_CATALOG_BASE_URL, data_val)
                if _is_catalog_product_link(normalized):
                    product_links.append(normalized)
        if anchor and hasattr(anchor, "get"):
            href = anchor.get("href")
            if href and isinstance(href, str):
                normalized = urljoin(_CATALOG_BASE_URL, href)
                if _is_catalog_product_link(normalized):
                    product_links.append(normalized)

    if not product_links:
        seen: set[str] = set()
        for a in page.select("a[href]"):
            href = a.get("href")
            if not isinstance(href, str):
                continue
            normalized = urljoin(_CATALOG_BASE_URL, href)
            if normalized in seen:
                continue
            if _is_catalog_product_link(normalized):
                seen.add(normalized)
                product_links.append(normalized)

    return product_links


def _detect_last_page(page: BeautifulSoup) -> int | None:
    last_page: int | None = None
    for anchor in page.select("ul.p52Mr a[href]"):
        href = anchor.get("href")
        if not isinstance(href, str):
            continue
        parsed = urlparse(urljoin(_CATALOG_BASE_URL, href))
        page_values = parse_qs(parsed.query).get("p")
        if not page_values:
            continue
        try:
            candidate = int(page_values[0])
        except (TypeError, ValueError):
            continue
        if last_page is None or candidate > last_page:
            last_page = candidate
    return last_page


def _is_catalog_product_link(link: str) -> bool:
    parsed = urlparse(link)
    if parsed.scheme not in {"http", "https"}:
        return False
    if not parsed.netloc.endswith("goldapple.ru"):
        return False
    return bool(
        _PRODUCT_URL_RE.match(f"{parsed.scheme}://{parsed.netloc}{parsed.path}")
    )


class GoldAppleScrapper(Scrapper):
    _pages: list[str]
    _product_batches: list[list[str]]

    def __init__(
        self,
        page_parser: PageParser,
        product_type_ids: Sequence[int] | Iterable[int] | None = None,
        *,
        max_pages: int | None = None,
        batch_size: int = 400,
    ):
        self._page_parser = page_parser
        self._product_type_ids = (
            list(product_type_ids)
            if product_type_ids is not None
            else list(_DEFAULT_PRODUCT_TYPE_IDS)
        )
        self._max_pages = max_pages
        self._batch_size = max(1, batch_size)

        self._product_links = _collect_product_links(
            self._product_type_ids,
            max_pages=self._max_pages,
        )
        print(f"Collected {len(self._product_links)} product links.")
        if not self._product_links:
            self._product_batches = []
            self._pages = []
        else:
            self._product_batches = [
                self._product_links[i : i + self._batch_size]
                for i in range(0, len(self._product_links), self._batch_size)
            ]
            self._pages = self._product_links

    def _is_product_link(self, link: str) -> bool:
        return _is_catalog_product_link(link)

    def scrap_page(self, index: int) -> list[Perfume]:
        if index + 1 > len(self._product_batches):
            return []

        batch_links = [
            link
            for link in self._product_batches[index]
            if link and self._is_product_link(link)
        ]
        if not batch_links:
            print(f"No product links to process for batch {index + 1}.")
            return []

        print(
            f"Scraping Gold Apple catalog batch {index + 1}/"
            f"{len(self._product_batches)} with {len(batch_links)} products."
        )

        perfumes = []
        locker = Lock()
        with tqdm(total=len(batch_links), desc="Scraping products") as pbar:
            with ThreadPoolExecutor(self._workers) as ex:
                futures = {
                    ex.submit(self.fetch_perfume, link): link for link in batch_links
                }
                for fut in as_completed(futures):
                    perfume = fut.result()
                    pbar.update(1)
                    if not perfume:
                        continue
                    with locker:
                        perfumes.append(perfume)

        print(f"Collected {len(perfumes)}.")
        return perfumes

    def fetch_perfume(self, link: str) -> Perfume | None:
        perfume_page = get_page(link, use_playwright=True)
        if not perfume_page:
            return None

        perfume = self._page_parser.parse_perfume_from_page(perfume_page)
        if not perfume:
            return None

        for volume_with_cost in perfume.shop_info.volumes_with_prices:
            volume_with_cost.link = link
        return perfume

    class PerfumeKey:
        brand: str
        name: str
        sex: str

        def __init__(self, perfume: Perfume):
            self.brand = perfume.brand
            self.name = perfume.name
            self.sex = perfume.sex

        def __hash__(self) -> int:
            return hash((self.brand, self.name, self.sex))

        def __eq__(self, other: object) -> bool:
            if not isinstance(other, self.__class__):
                return False
            return (
                self.brand == other.brand
                and self.name == other.name
                and self.sex == other.sex
            )

    def scrap_all_accuratly(self) -> list[Perfume]:
        perfumes = []
        for i in range(len(self._product_batches)):
            page_perfumes = self.scrap_page(i)
            perfumes.extend(page_perfumes)

        perfumes_with_glued_links: dict[GoldAppleScrapper.PerfumeKey, Perfume] = {}
        for perfume in perfumes:
            key = self.PerfumeKey(perfume)
            if key in perfumes_with_glued_links:
                perfumes_with_glued_links[key].shop_info.volumes_with_prices.extend(
                    perfume.shop_info.volumes_with_prices
                )
            else:
                perfumes_with_glued_links[key] = perfume

        return list(perfumes_with_glued_links.values())
