from concurrent.futures import ThreadPoolExecutor, as_completed
from threading import Lock

from tqdm import tqdm

from models import Perfume
from scraping import PageParser, Scrapper
from util import get_page


class RandewooScrapper(Scrapper):
    def __init__(self, page_parser: PageParser):
        self._page_parser = page_parser
        self._perfume_catalog_link = (
            "https://randewoo.ru/category/parfyumeriya?paging=200"
        )

    def scrap_page(self, index: int) -> list[Perfume]:
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

        perfumes = []
        locker = Lock()
        with tqdm(total=len(perfume_links), desc=f"Scraping page {index + 1}") as pbar:
            with ThreadPoolExecutor(self._workers) as ex:
                futures = {
                    ex.submit(self.fetch_perfume, link): link for link in perfume_links
                }
                for fut in as_completed(futures):
                    perfume = fut.result()
                    pbar.update(1)
                    if not perfume:
                        continue
                    with locker:
                        perfumes.append(perfume)

        print(f"Collected {len(perfumes)} perfumes from page {index + 1}.")
        return perfumes

    def fetch_perfume(self, link: str) -> Perfume | None:
        perfume_page = get_page(link)
        if not perfume_page:
            return None

        perfume = self._page_parser.parse_perfume_from_page(perfume_page)
        if not perfume:
            return None

        for volume_with_cost in perfume.shop_info.volumes_with_prices:
            volume_with_cost.link = link
        return perfume

    def _normalize_link(self, link: str) -> str:
        if link.startswith("/"):
            return "https://randewoo.ru" + link
        return link

    def scrap_all_accuratly(self) -> list[Perfume]:
        perfumes = []
        i = 0
        while True:
            page_perfumes = self.scrap_page(i)
            if not page_perfumes:
                break
            perfumes.extend(page_perfumes)
            i += 1
        return perfumes
