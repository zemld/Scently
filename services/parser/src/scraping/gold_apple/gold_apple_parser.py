import re

from bs4 import BeautifulSoup, Tag

from src.models import PerfumeFromConcreteShop

from ..page_parser import PageParser


class GoldApplePageParser(PageParser):
    _split_notes_pattern = r",\s*|\s+и\s+|\s+-\s+|\s+–\s+|\s+"

    def _parse_brand(self, page: BeautifulSoup) -> str:
        brand_tag = page.find("a", {"data-transaction-name": "ga-pdp-title"})
        if brand_tag:
            brand = brand_tag.get_text(strip=True)
            if brand:
                return str(brand)

        h1_tag = page.find("h1")
        if h1_tag:
            brand_link = h1_tag.find("a", {"content": True})
            if isinstance(brand_link, Tag):
                content = brand_link.get("content")
                if content and isinstance(content, str):
                    return content
            brand_link = h1_tag.find("a")
            if isinstance(brand_link, Tag):
                return str(brand_link.get_text(strip=True))

        meta_keywords = page.find("meta", {"name": "keywords"})
        if isinstance(meta_keywords, Tag):
            content = meta_keywords.get("content")
            if content and isinstance(content, str):
                keywords = content
                words = keywords.split()
                if len(words) >= 3:
                    return words[2]

        return ""

    def _parse_name(self, page: BeautifulSoup) -> str:
        h1_tag = page.find("h1")
        if h1_tag:
            name_span = h1_tag.find("span")
            if isinstance(name_span, Tag):
                name = name_span.get_text(strip=True)
                if name:
                    return str(name)
        return ""

    def _parse_props(self, page: BeautifulSoup) -> dict[str, str]:
        props: dict[str, str] = {}

        dl_elements = page.find_all("dl")
        for dl in dl_elements:
            direct_divs = dl.find_all("div", recursive=False)
            if direct_divs:
                first_div = direct_divs[0]
                prop_divs = first_div.find_all("div", recursive=False)
            else:
                prop_divs = dl.find_all("div", recursive=True)

            for prop_div in prop_divs:
                dt_elements = prop_div.find_all("dt", limit=2)
                if len(dt_elements) < 2:
                    continue
                key_span = dt_elements[0].find("span")
                value_span = dt_elements[1].find("span")
                if not key_span or not value_span:
                    continue
                key = key_span.get_text(strip=True).lower()
                value = value_span.get_text(strip=True)
                if not key or not value:
                    continue
                props[key] = value

        return props

    def _parse_type(self, page: BeautifulSoup) -> str:
        props = self._parse_props(page)
        type_from_props = props.get("тип продукта", "")
        if type_from_props:
            type_text = type_from_props.lower().strip()
            if not self._type_canonizer:
                return type_text
            canonized_type = self._type_canonizer.canonize(type_text)
            return canonized_type if canonized_type else type_text

        meta_desc = page.find("meta", {"name": "description"})
        if isinstance(meta_desc, Tag):
            content = meta_desc.get("content")
            if content and isinstance(content, str):
                desc = content.lower()
                if not self._type_canonizer:
                    return desc
                canonized_type = self._type_canonizer.canonize(desc)
                if canonized_type:
                    return canonized_type

        title_tag = page.find("title")
        if title_tag:
            title = title_tag.get_text(strip=True).lower()
            if not self._type_canonizer:
                return title
            canonized_type = self._type_canonizer.canonize(title)
            if canonized_type:
                return canonized_type
        return ""

    def _parse_sex(self, props: dict[str, str]) -> str:
        raw_sex = props.get("для кого")
        if not raw_sex:
            return ""
        if " " in raw_sex:
            raw_sex = raw_sex.split()[-1]
        if not self._sex_canonizer:
            return raw_sex
        canonized_sex = self._sex_canonizer.canonize(raw_sex)
        if canonized_sex:
            return canonized_sex
        return raw_sex

    def _parse_families(self, props: dict[str, str]) -> list[str]:
        families = props.get("группа ароматов")
        if not families:
            return []
        if not self._family_canonizer:
            return [family.strip() for family in families.split(",") if family.strip()]
        canonized_families = [
            self._family_canonizer.canonize(family.strip())
            for family in families.split(",")
        ]
        unique_families = list(set(canonized_families))
        return [family for family in unique_families if family]

    def _parse_notes(self, props: dict[str, str], key: str) -> list[str]:
        notes = props.get(key)
        if not notes:
            return []
        notes_list = re.split(self._split_notes_pattern, notes)
        if not self._notes_canonizer:
            return [note.strip() for note in notes_list if note.strip()]
        canonized_notes = [
            self._notes_canonizer.canonize(note.strip()) for note in notes_list
        ]
        return [note for note in canonized_notes if note]

    def _parse_upper_notes(self, props: dict[str, str]) -> list[str]:
        return self._parse_notes(props, "верхние ноты")

    def _parse_core_notes(self, props: dict[str, str]) -> list[str]:
        return self._parse_notes(props, "средние ноты")

    def _parse_base_notes(self, props: dict[str, str]) -> list[str]:
        return self._parse_notes(props, "базовые ноты")

    def _get_shop_info(self, page: BeautifulSoup) -> PerfumeFromConcreteShop.ShopInfo:
        shop_info = PerfumeFromConcreteShop.ShopInfo(
            shop_name="Gold Apple",
            domain="https://goldapple.ru",
            image_url=self._parse_image_url(page),
            variants=[],
        )
        current_item_variant = self._parse_current_item_variant(page)
        if current_item_variant:
            shop_info.variants.append(current_item_variant)
        return shop_info

    def _extract_volume(self, page: BeautifulSoup) -> int | None:
        props = self._parse_props(page)
        if not props:
            return None
        volume = props.get("объем") or props.get("объём")
        if not volume:
            return None

        volume_text = volume.split()[0]
        try:
            volume_value = int(volume_text)
        except ValueError:
            return None
        return int(volume_value)

    def _extract_price(self, page: BeautifulSoup) -> int | None:
        price_meta = page.find("meta", itemprop="price")
        if isinstance(price_meta, Tag):
            content = price_meta.get("content")
            if content and isinstance(content, str):
                try:
                    return int(content)
                except (ValueError, TypeError):
                    pass

        price_element = page.find(itemprop="price")
        if price_element:
            price_text = price_element.get_text(strip=True)
            price_digits = "".join(c for c in price_text if c.isdigit())
            if price_digits:
                try:
                    return int(price_digits)
                except ValueError:
                    pass
        return None

    def _extract_link(self, page: BeautifulSoup) -> str:
        link_tag = page.find("link", itemprop="url")
        if not isinstance(link_tag, Tag):
            return ""
        link_text = link_tag.get("href")
        if not isinstance(link_text, str):
            return ""
        return link_text

    def _normalize_link(self, domain: str, link: str) -> str:
        if link.startswith("/"):
            return domain + link
        return link

    def _parse_current_item_variant(
        self, page: BeautifulSoup
    ) -> PerfumeFromConcreteShop.ShopInfo.VolumeWithPrices | None:
        volume = self._extract_volume(page)
        if not volume:
            return None
        price = self._extract_price(page)
        if not price:
            return None
        link = self._extract_link(page)
        if not link:
            return None
        return PerfumeFromConcreteShop.ShopInfo.VolumeWithPrices(volume, price, link)

    def _parse_image_url(self, page: BeautifulSoup) -> str:
        og_image = page.find("meta", property="og:image")
        if isinstance(og_image, Tag):
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
