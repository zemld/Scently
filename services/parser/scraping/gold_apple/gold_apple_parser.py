import re

from bs4 import BeautifulSoup, element

from scraping.page_parser import PageParser


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
        if self._brand_canonizer:
            result = super()._canonize(brand_info[0], self._brand_canonizer)
            if isinstance(result, str):
                return result
            return str(brand_info[0])
        return str(brand_info[0])

    @staticmethod
    def _is_brand_tag(tag: element.Tag) -> bool:
        return tag.has_attr("text") and tag.get("text") == "Бренд"

    def _parse_name(self, page: BeautifulSoup) -> str:
        name_tag = page.find_all(self._is_name_tag)
        if not name_tag or not name_tag[0].string:
            return ""
        if self._name_canonizer:
            result = super()._canonize(name_tag[0].string.strip(), self._name_canonizer)
            if isinstance(result, str):
                return result
            return str(name_tag[0].string.strip())
        return str(name_tag[0].string.strip())

    @staticmethod
    def _is_name_tag(tag: element.Tag) -> bool:
        return (
            tag.name == "span"
            and tag.has_attr("itemprop")
            and tag.get("itemprop") == "name"
            and tag.has_attr("class")
        )

    def _parse_properties(self, page: BeautifulSoup) -> list[str]:
        properties_title_rx = re.compile("Подробные характеристики", re.IGNORECASE)
        properties_title = page.find_all(string=properties_title_rx)
        if not properties_title:
            return []
        try:
            if not properties_title or not properties_title[0].parent:
                return []
            section = properties_title[0].parent.parent
            if section is None:
                return []
            raw_properties = section.find_all("span")
            properties = []
            for prop in raw_properties:
                if prop.string is not None:
                    properties.append(prop.string.strip().lower())
            return properties
        except Exception:
            return []

    def _parse_type(self, page: BeautifulSoup) -> str:
        props = self._parse_properties(page)
        if not props or len(props) < 2:
            return ""
        if self._type_canonizer:
            result = super()._canonize(props[1], self._type_canonizer)
            return result if isinstance(result, str) else props[1]
        return props[1]

    def _parse_sex(self, page: BeautifulSoup) -> str:
        props = self._parse_properties(page)
        if not props or len(props) < 4:
            return ""
        if self._sex_canonizer:
            result = super()._canonize(props[3].lower().split(), self._sex_canonizer)
            return result if isinstance(result, str) else props[3]
        return props[3]

    def _parse_families(self, page: BeautifulSoup) -> list[str]:
        props = self._parse_properties(page)
        if not props or len(props) < 6:
            return []
        families = []
        for family in props[5].split(","):
            if self._family_canonizer:
                result = super()._canonize(
                    family.strip().lower(), self._family_canonizer
                )
                families.append(
                    result if isinstance(result, str) else family.strip().lower()
                )
            else:
                families.append(family.strip().lower())
        return [family for family in families if family]

    def _parse_notes(self, notes: str) -> list[str]:
        notes_list = [
            note.strip(" .,").lower()
            for note in re.split(self._split_notes_pattern, notes)
            if note.strip()
        ]
        if self._notes_canonizer:
            canonized = super()._canonize(notes_list, self._notes_canonizer)
            return [note for note in canonized if note]
        return [note for note in notes_list if note]

    def _parse_upper_notes(self, page: BeautifulSoup) -> list[str]:
        props = self._parse_properties(page)
        if not props or len(props) < 8:
            return []
        return self._parse_notes(props[7])

    def _parse_middle_notes(self, page: BeautifulSoup) -> list[str]:
        props = self._parse_properties(page)
        if not props or len(props) < 10:
            return []
        return self._parse_notes(props[9])

    def _parse_base_notes(self, page: BeautifulSoup) -> list[str]:
        props = self._parse_properties(page)
        if not props or len(props) < 12:
            return []
        return self._parse_notes(props[11])

    def _parse_volume(self, page: BeautifulSoup) -> int:
        props = self._parse_properties(page)
        if not props or len(props) < 14:
            return 0
        int_rx = re.compile(r"\d+")
        match = re.search(int_rx, props[13])
        return int(match.group(0)) if match else 0

    def _parse_image_url(self, page: BeautifulSoup) -> str:
        og_image = page.find("meta", property="og:image")
        if og_image and hasattr(og_image, "get") and og_image.get("content"):
            image_url = og_image.get("content")
            if isinstance(image_url, str) and self._is_valid_image_url(image_url):
                return image_url

        return ""

    @staticmethod
    def _is_valid_image_url(url: str) -> bool:
        if not url or not isinstance(url, str):
            return False

        invalid_patterns = [
            "placeholder",
            "no-image",
            "default",
            "empty",
            "blank",
            "loading",
            "spinner",
        ]

        url_lower = url.lower()
        if any(pattern in url_lower for pattern in invalid_patterns):
            return False

        valid_extensions = [".jpg", ".jpeg", ".png", ".webp", ".gif"]
        return any(url_lower.endswith(ext) for ext in valid_extensions) or "?" in url
