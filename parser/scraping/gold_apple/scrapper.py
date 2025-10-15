from scraping.scrapper import Scrapper
from util.send_request import get_page
from models.perfume import Perfume
import re
from concurrent.futures import ThreadPoolExecutor, as_completed
from threading import Lock
import time
from scraping.page_parser import PageParser
from tqdm import tqdm


class GoldAppleScrapper(Scrapper):
    _product_url_re = re.compile(
        r"^https://goldapple\.ru/\d{6,}-[a-z0-9-]+$", re.IGNORECASE
    )

    def __init__(self, page_parser: PageParser):
        sitemaps_url = "https://goldapple.ru/sitemap.xml"

        sitemaps = [
            sitemap.find("loc").string
            for sitemap in get_page(sitemaps_url).find_all("sitemap")
        ]
        self._sitemaps = [sitemap for sitemap in sitemaps if sitemap]
        self._page_parser = page_parser

    def _is_product_link(self, link: str) -> bool:
        return bool(self._product_url_re.match(link))

    def scrap_sitemap(self, index) -> list[Perfume]:
        print(f"Scraping sitemap {self._sitemaps[index]}.")
        if index + 1 > len(self._sitemaps):
            return None
        links_page = get_page(self._sitemaps[index])
        if not links_page:
            print("Unable to scrap.")
            return []

        links = [link.string.strip() for link in links_page.find_all("loc")]
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

    def fetch_perfume(self, link) -> Perfume | None:
        time.sleep(1)

        page = get_page(link)
        if not page:
            return None

        if not any(rx.search(page.title.string.strip()) for rx in self._perfumes_re):
            return None

        perfume = self._page_parser.parse_perfume_from_page(page)
        if not perfume:
            return None

        perfume.link = link
        return perfume
