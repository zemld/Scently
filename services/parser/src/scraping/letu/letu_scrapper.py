from src.models import Perfume
from src.util import get_page

from ..page_parser import PageParser
from ..scrapper import Scrapper


class LetuScrapper(Scrapper):
    def __init__(self, domain: str, page_parser: PageParser):
        self._domain = domain
        self._page_parser = page_parser
        self._perfume_catalog_link = "https://www.letu.ru/browse/parfyumeriya/filters/product-class=duhi-or-parfyumernaya-voda-or-tualetnaya-voda"
        self._workers = 8

    def scrap_page(self, index: int) -> list[Perfume]:
        page_url = self._perfume_catalog_link
        if index != 0:
            page_url = f"{self._perfume_catalog_link}/page={index + 1}"
        page = get_page(page_url, use_playwright=True)
        if not page:
            print(f"Failed to load catalog page {page_url}")
            return []

        perfume_link_tags = page.find_all("a", class_="product-tile__item-container")
        perfume_links: list[str] = []
        for tag in perfume_link_tags:
            href = tag.get("href")
            if not isinstance(href, str):
                continue
            perfume_links.append(self._normalize_link(href))

        print(f"Found {len(perfume_links)} links on catalog page {page_url}")
        return self.process_page_links(perfume_links, index)

    def fetch_perfume(self, link: str) -> Perfume | None:
        perfume_page = get_page(link, use_playwright=True)
        if not perfume_page:
            print(f"Failed to load perfume page {link}")
            return None
        perfume = self._page_parser.parse_perfume_from_page(perfume_page)
        if not perfume:
            print(f"Failed to parse perfume page {link}")
            return None
        return perfume
