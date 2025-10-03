from scraping.page_parser import PageParser
from bs4 import BeautifulSoup, element
import re


class GoldApplePageParser(PageParser):
    _split_notes_pattern = r",\s*|\s+и\s+|\s+-\s+|\s+–\s+"

    def _parse_brand(self, page: BeautifulSoup) -> str:
        brand_tags = page.find_all(self._is_brand_tag)
        if not brand_tags:
            return ""

        brand_info = [
            tag.string.strip() for tag in brand_tags[0].find_all("div") if tag.string
        ]
        if not brand_info:
            return ""
        return super()._canonize(brand_info[0], self._brand_canonizer)

    def _is_brand_tag(tag: element.Tag) -> bool:
        return tag.has_attr("text") and tag.get("text") == "Бренд"

    def _parse_name(self, page: BeautifulSoup) -> str:
        name_tag = page.find_all(self._is_name_tag)
        if not name_tag or not name_tag[0].string:
            return ""
        return super()._canonize(name_tag[0].string.strip(), self._name_canonizer)

    def _is_name_tag(tag: element.Tag) -> bool:
        return (
            tag.name == "span"
            and tag.has_attr("itemprop")
            and tag.get("itemprop") == "name"
            and tag.has_attr("class")
        )

    def _parse_properties(page: BeautifulSoup) -> list[str]:
        properties_title_rx = re.compile("Подробные характеристики", re.IGNORECASE)
        properties_title = page.find_all(string=properties_title_rx)
        if not properties_title:
            return []
        try:
            section = properties_title[0].parent.parent if properties_title else None
            raw_properties = section.find_all("span")
            properties = []
            for prop in raw_properties:
                properties.append(prop.string.strip().lower())
            return properties
        except Exception as e:
            return []

    def _parse_type(self, page) -> str:
        props = self._parse_properties(page)
        if not props or len(props) < 2:
            return ""
        return super()._canonize(props[1], self._type_canonizer)

    def _parse_sex(self, page) -> str:
        props = self._parse_properties(page)
        if not props or len(props) < 4:
            return ""
        return super()._canonize(props[3], self._sex_canonizer)

    def _parse_families(self, page) -> list[str]:
        props = self._parse_properties(page)
        if not props or len(props) < 6:
            return []
        families = [
            super()._canonize(family.strip().lower(), self._family_canonizer)
            for family in props[5].split(",")
        ]
        return [family for family in families if family]

    def _parse_notes(self, notes: str) -> list[str]:
        notes_list = [
            super()._canonize(note.strip().lower(), self._notes_canonizer)
            for note in re.split(self._split_notes_pattern, notes)
            if note.strip()
        ]
        return [note for note in notes_list if note]

    def _parse_upper_notes(self, page) -> list[str]:
        props = self._parse_properties(page)
        if not props or len(props) < 8:
            return []
        return self._parse_notes(props[7])

    def _parse_middle_notes(self, page) -> list[str]:
        props = self._parse_properties(page)
        if not props or len(props) < 10:
            return []
        return self._parse_notes(props[9])

    def _parse_base_notes(self, page) -> list[str]:
        props = self._parse_properties(page)
        if not props or len(props) < 12:
            return []
        return self._parse_notes(props[11])

    def _parse_volume(self, page) -> int:
        props = self._parse_properties(page)
        if not props or len(props) < 14:
            return 0
        int_rx = re.compile(r"\d+")
        return (
            int(re.search(int_rx, props[13]).group(0))
            if re.search(int_rx, props[13])
            else 0
        )
