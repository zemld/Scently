import re
from abc import ABC, abstractmethod

from models.perfume import Perfume
from scraping.page_parser import PageParser


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

    @abstractmethod
    def scrap_page(self, index: int) -> list[Perfume]:
        pass

    @abstractmethod
    def fetch_perfume(self, link: str) -> Perfume | None:
        pass

    @abstractmethod
    def scrap_all_accuratly(self) -> list[Perfume]:
        pass
