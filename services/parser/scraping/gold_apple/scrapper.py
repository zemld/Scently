import re
import time
from concurrent.futures import ThreadPoolExecutor, as_completed
from threading import Lock

from tqdm import tqdm

from models.perfume import Perfume
from scraping.page_parser import PageParser
from scraping.scrapper import Scrapper
from util.send_request import get_page


class GoldAppleScrapper(Scrapper):
    _product_url_re = re.compile(
        r"^https://goldapple\.ru/\d{6,}-[a-z0-9-]+$", re.IGNORECASE
    )

    def __init__(self, page_parser: PageParser):
        sitemaps_url = "https://goldapple.ru/sitemap.xml"

        sitemap_page = get_page(sitemaps_url)
        sitemaps: list[str] = []
        if sitemap_page:
            for sitemap in sitemap_page.find_all("sitemap"):
                loc_tag = sitemap.find("loc")
                if loc_tag and loc_tag.string:
                    sitemaps.append(loc_tag.string)
        self._pages = [sitemap for sitemap in sitemaps if sitemap]
        self._page_parser = page_parser

    def _is_product_link(self, link: str) -> bool:
        return bool(self._product_url_re.match(link))

    def scrap_page(self, index: int) -> list[Perfume]:
        print(f"Scraping sitemap {self._pages[index]}.")
        if index + 1 > len(self._pages):
            return []
        links_page = get_page(self._pages[index])
        if not links_page:
            print("Unable to scrap.")
            return []

        links = []
        for link in links_page.find_all("loc"):
            if link.string:
                links.append(link.string.strip())
        product_links = [link for link in links if link and self._is_product_link(link)]
        print(f"Found {len(product_links)} in sitemap.")

        perfumes = []
        locker = Lock()
        with tqdm(total=len(product_links), desc="Scraping products") as pbar:
            with ThreadPoolExecutor(self._workers) as ex:
                futures = {
                    ex.submit(self.fetch_perfume, link): link for link in product_links
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
        time.sleep(1)

        page = get_page(link)
        if not page:
            return None

        if (
            not page.title
            or not page.title.string
            or not any(rx.search(page.title.string.strip()) for rx in self._perfumes_re)
        ):
            return None

        perfume = self._page_parser.parse_perfume_from_page(page)
        if not perfume:
            return None

        perfume.link = link
        return perfume
