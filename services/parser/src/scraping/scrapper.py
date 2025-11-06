import re
from abc import ABC, abstractmethod
from concurrent.futures import ThreadPoolExecutor, as_completed
from threading import Lock

from tqdm import tqdm

from src.models.perfume import Perfume
from src.scraping.page_parser import PageParser


class Scrapper(ABC):
    _pages: list[str]
    _workers: int = 16
    _perfumes_re = [
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
    _page_parser: PageParser
    _domain: str

    @abstractmethod
    def scrap_page(self, index: int) -> list[Perfume]:
        pass

    def process_page_links(
        self, page_links: list[str], page_index: int
    ) -> list[Perfume]:
        perfumes = []
        locker = Lock()
        with tqdm(
            total=len(page_links), desc=f"Scraping page {page_index + 1}", leave=False
        ) as pbar:
            with ThreadPoolExecutor(self._workers) as ex:
                futures = {
                    ex.submit(self.fetch_perfume, link): link for link in page_links
                }
                for fut in as_completed(futures):
                    perfume = fut.result()
                    pbar.update(1)
                    if not perfume:
                        continue
                    with locker:
                        perfumes.append(perfume)
        print(f"Collected {len(perfumes)} perfumes from page {page_index + 1}.")
        return perfumes

    @abstractmethod
    def fetch_perfume(self, link: str) -> Perfume | None:
        pass

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

    def _normalize_link(self, link: str) -> str:
        if link.startswith("/"):
            return self._domain + link
        return link
