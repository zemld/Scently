from abc import ABC, abstractmethod

from bs4 import BeautifulSoup

from src.canonization.canonizer import Canonizer
from src.models.perfume import PerfumeFromConcreteShop


class PageParser(ABC):
    _type_canonizer: Canonizer | None
    _sex_canonizer: Canonizer | None
    _family_canonizer: Canonizer | None
    _notes_canonizer: Canonizer | None

    def __init__(
        self,
        type_canonizer: Canonizer | None = None,
        sex_canonizer: Canonizer | None = None,
        family_canonizer: Canonizer | None = None,
        notes_canonizer: Canonizer | None = None,
    ):
        self._type_canonizer = type_canonizer
        self._sex_canonizer = sex_canonizer
        self._family_canonizer = family_canonizer
        self._notes_canonizer = notes_canonizer

    def parse_perfume_from_page(
        self, page: BeautifulSoup
    ) -> PerfumeFromConcreteShop | None:
        brand = self._parse_brand(page)
        name = self._parse_name(page)
        perfume_type = self._parse_type(page)

        props = self._parse_props(page)
        sex = self._parse_sex(props)
        families = self._parse_families(props)
        upper_notes = self._parse_upper_notes(props)
        core_notes = self._parse_core_notes(props)
        base_notes = self._parse_base_notes(props)

        shop_info = self._get_shop_info(page)
        if any(
            not item for item in (brand, name, perfume_type, sex, shop_info.image_url)
        ):
            return None

        perfume_properties = PerfumeFromConcreteShop.PerfumeProperties(
            perfume_type=perfume_type,
            family=families,
            upper_notes=upper_notes,
            core_notes=core_notes,
            base_notes=base_notes,
        )

        perfume = PerfumeFromConcreteShop(
            brand=brand,
            name=name,
            sex=sex,
            properties=perfume_properties,
            shop_info=shop_info,
        )
        return perfume

    @abstractmethod
    def _parse_brand(self, page: BeautifulSoup) -> str:
        pass

    @abstractmethod
    def _parse_name(self, page: BeautifulSoup) -> str:
        pass

    @abstractmethod
    def _parse_props(self, page: BeautifulSoup) -> dict[str, str]:
        pass

    @abstractmethod
    def _parse_type(self, page: BeautifulSoup) -> str:
        pass

    @abstractmethod
    def _parse_sex(self, props: dict[str, str]) -> str:
        pass

    @abstractmethod
    def _parse_families(self, props: dict[str, str]) -> list[str]:
        pass

    @abstractmethod
    def _parse_notes(self, props: dict[str, str], key: str) -> list[str]:
        pass

    @abstractmethod
    def _parse_upper_notes(self, props: dict[str, str]) -> list[str]:
        pass

    @abstractmethod
    def _parse_core_notes(self, props: dict[str, str]) -> list[str]:
        pass

    @abstractmethod
    def _parse_base_notes(self, props: dict[str, str]) -> list[str]:
        pass

    @abstractmethod
    def _get_shop_info(self, page: BeautifulSoup) -> PerfumeFromConcreteShop.ShopInfo:
        pass

    @abstractmethod
    def _parse_image_url(self, page: BeautifulSoup) -> str:
        pass
