from abc import ABC, abstractmethod
from models.perfume import Perfume
import time


class Scrapper(ABC):
    sitemaps: list[str]

    @abstractmethod
    def scrap_sitemap(self, index: int) -> list[Perfume]:
        pass

    def scrap_all_accuratly(self) -> list[Perfume]:
        perfumes = []
        for i in range(len(self.sitemaps)):
            perfumes.append(self.scrap_sitemap(i))
            time.sleep(3600)
        return perfumes
