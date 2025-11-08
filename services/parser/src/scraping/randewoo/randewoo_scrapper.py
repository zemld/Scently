from pathlib import Path

from src.models import PerfumeFromConcreteShop
from src.util import get_page
from src.util.backup import BackupManager

from ..page_parser import PageParser
from ..scrapper import Scrapper


class RandewooScrapper(Scrapper):
    def __init__(
        self,
        domain: str,
        page_parser: PageParser,
        backup_dir: Path | None = None,
    ):
        self._domain = domain
        self._page_parser = page_parser
        self._perfume_catalog_link = (
            "https://randewoo.ru/category/parfyumeriya?paging=200"
        )
        self._backup_manager = BackupManager("randewoo", backup_dir)

    def scrap_page(self, index: int) -> list[PerfumeFromConcreteShop]:
        page_url = f"{self._perfume_catalog_link}&page={index + 1}"
        page = get_page(page_url)
        if not page:
            return []

        perfume_link_tags = page.find_all("a", class_="b-catalogItem__photoWrap")
        perfume_links: list[str] = []
        for tag in perfume_link_tags:
            href = tag.get("href")
            if not isinstance(href, str):
                continue
            perfume_links.append(self._normalize_link(href))

        existing_links = self._backup_manager.load_links()
        new_links = [link for link in perfume_links if link not in existing_links]

        if new_links:
            self._backup_manager.add_links(new_links)
            print(f"Saved {len(new_links)} new links to backup.")

        return self.process_page_links(perfume_links, index, self._backup_manager)

    def fetch_perfume(self, link: str) -> PerfumeFromConcreteShop | None:
        perfume_page = get_page(link)
        if not perfume_page:
            return None

        perfume = self._page_parser.parse_perfume_from_page(perfume_page)
        if not perfume:
            return None

        for volume_with_cost in perfume.shop_info.variants:
            volume_with_cost.link = link
        return perfume
