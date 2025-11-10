import re
from abc import ABC, abstractmethod
from concurrent.futures import ThreadPoolExecutor, as_completed
from pathlib import Path
from threading import Lock

from tqdm import tqdm

from src.models.perfume import PerfumeFromConcreteShop
from src.scraping.page_parser import PageParser
from src.util import BackupManager, setup_logger

scrapper_logger = setup_logger(
    __name__, log_file=Path.cwd() / "logs" / f"{__name__.split('.')[-1]}.log"
)


class Scrapper(ABC):
    _pages: list[str]
    _workers: int = 8
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
    def scrap_page(self, index: int) -> list[PerfumeFromConcreteShop]:
        pass

    def process_page_links(
        self,
        page_links: list[str],
        page_index: int,
        backup_manager: BackupManager | None = None,
    ) -> list[PerfumeFromConcreteShop]:
        perfumes = []
        locker = Lock()
        with tqdm(
            total=len(page_links), desc=f"Scraping page {page_index + 1}", leave=False
        ) as pbar:
            with ThreadPoolExecutor(self._workers) as ex:
                futures = {
                    ex.submit(self.fetch_perfume, link): link for link in page_links
                }
                batch_perfumes = []
                for fut in as_completed(futures):
                    perfume = fut.result()
                    pbar.update(1)
                    if not perfume:
                        continue
                    with locker:
                        perfumes.append(perfume)
                        batch_perfumes.append(perfume)

                if backup_manager and batch_perfumes:
                    backup_manager.add_perfumes(batch_perfumes)

        scrapper_logger.info(
            f"Collected perfumes from page | page_index={page_index} | "
            f"count={len(perfumes)}"
        )
        return perfumes

    @abstractmethod
    def fetch_perfume(self, link: str) -> PerfumeFromConcreteShop | None:
        pass

    def scrap_all_accuratly(self) -> list[PerfumeFromConcreteShop]:
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
