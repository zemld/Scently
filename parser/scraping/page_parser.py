from abc import ABC, abstractmethod
from bs4 import BeautifulSoup
from models.perfume import Perfume
from canonization.canonizer import Canonizer


class PageParser(ABC):
    _brand_canonizer: Canonizer
    _name_canonizer: Canonizer
    _type_canonizer: Canonizer
    _sex_canonizer: Canonizer
    _family_canonizer: Canonizer
    _notes_canonizer: Canonizer

    def __init__(
        self,
        brand_canonizer: Canonizer = None,
        name_canonizer: Canonizer = None,
        type_canonizer: Canonizer = None,
        sex_canonizer: Canonizer = None,
        family_canonizer: Canonizer = None,
        notes_canonizer: Canonizer = None,
    ):
        self._brand_canonizer = brand_canonizer
        self._name_canonizer = name_canonizer
        self._type_canonizer = type_canonizer
        self._sex_canonizer = sex_canonizer
        self._family_canonizer = family_canonizer
        self._notes_canonizer = notes_canonizer

    def parse_perfume_from_page(self, page: BeautifulSoup) -> Perfume | None:
        brand = self._parse_brand(page)
        name = self._parse_name(page)
        perfume_type = self._parse_type(page)
        sex = self._parse_sex(page)
        families = self._parse_families(page)
        upper_notes = self._parse_upper_notes(page)
        middle_notes = self._parse_middle_notes(page)
        base_notes = self._parse_base_notes(page)
        volume = self._parse_volume(page)

        if (
            any(item == "" for item in (brand, name, perfume_type, sex))
            or any(
                len(item) for item in (families, upper_notes, middle_notes, base_notes)
            )
            or volume == 0
        ):
            return None

        return Perfume(
            brand,
            name,
            perfume_type,
            sex,
            families,
            upper_notes,
            middle_notes,
            base_notes,
            volume,
        )

    @abstractmethod
    def _parse_brand(self, page: BeautifulSoup) -> str:
        pass

    @abstractmethod
    def _parse_name(self, page: BeautifulSoup) -> str:
        pass

    @abstractmethod
    def _parse_type(self, page: BeautifulSoup) -> str:
        pass

    @abstractmethod
    def _parse_sex(self, page: BeautifulSoup) -> str:
        pass

    @abstractmethod
    def _parse_families(self, page: BeautifulSoup) -> list[str]:
        pass

    @abstractmethod
    def _parse_upper_notes(self, page: BeautifulSoup) -> list[str]:
        pass

    @abstractmethod
    def _parse_middle_notes(self, page: BeautifulSoup) -> list[str]:
        pass

    @abstractmethod
    def _parse_base_notes(self, page: BeautifulSoup) -> list[str]:
        pass

    @abstractmethod
    def _parse_volume(self, page: BeautifulSoup) -> int:
        pass
