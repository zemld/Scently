import re

from bs4 import BeautifulSoup, Tag

from src.models import PerfumeFromConcreteShop
from src.util import get_page

from ..page_parser import PageParser


class LetuPageParser(PageParser):
    def _parse_brand(self, page: BeautifulSoup) -> str:
        brand_tag = page.find("meta", itemprop="name")
        if brand_tag and isinstance(brand_tag, Tag):
            content = brand_tag.get("content")
            if isinstance(content, str):
                return content
        return ""

    def _parse_header(self, page: BeautifulSoup) -> str:
        header_tag = page.find("h1", class_="product-detail-v3__name-title")
        if not header_tag:
            return ""
        return header_tag.get_text(strip=True)

    def _parse_name(self, page: BeautifulSoup) -> str:
        header = self._parse_header(page)
        if not header:
            return ""
        brand = self._parse_brand(page)
        return header.split(",")[0].strip()[len(brand) :].strip()

    def _parse_type(self, page: BeautifulSoup) -> str:
        header = self._parse_header(page)
        header_parts = header.split(",")
        if len(header_parts) < 2:
            return ""
        if not self._type_canonizer:
            return header_parts[1].strip()
        return self._type_canonizer.canonize(header_parts[1].strip())

    def _parse_props(self, page: BeautifulSoup) -> dict[str, str]:
        props: dict[str, str] = {}
        prop_tags = page.find_all("div", class_="product-group-characteristics__item")
        for prop_tag in prop_tags:
            prop_name_tag = prop_tag.find(
                "div", class_="product-group-characteristics__item-name"
            )
            prop_value_tag = prop_tag.find(
                "span", class_="product-group-characteristics__item-value"
            )
            if not prop_name_tag or not prop_value_tag:
                continue
            prop_name = prop_name_tag.get_text(strip=True)
            prop_value = prop_value_tag.get_text(strip=True)
            props[prop_name.lower()] = prop_value.lower()
        return props

    def _parse_sex(self, props: dict[str, str]) -> str:
        raw_sex = props.get("пол")
        if not raw_sex:
            return ""
        if not self._sex_canonizer:
            return raw_sex
        return self._sex_canonizer.canonize(raw_sex)

    def _parse_families(self, props: dict[str, str]) -> list[str]:
        families = props.get("группа аромата")
        if not families:
            return []

        splittered_families = families.split(",")
        if not self._family_canonizer:
            return [family.strip() for family in splittered_families if family.strip()]
        canonized_families = [
            self._family_canonizer.canonize(family.strip())
            for family in splittered_families
        ]
        unique_families = list(set(canonized_families))
        return [family for family in unique_families if family]

    def _parse_notes(self, props: dict[str, str], key: str) -> list[str]:
        notes = props.get(key)
        if not notes:
            return []

        notes_list = notes.split(",")
        if not self._notes_canonizer:
            return [note.strip() for note in notes_list if note.strip()]
        canonized_notes = [
            self._notes_canonizer.canonize(note.strip()) for note in notes_list
        ]
        return [note for note in canonized_notes if note]

    def _parse_upper_notes(self, props: dict[str, str]) -> list[str]:
        return self._parse_notes(props, "верхние ноты")

    def _parse_core_notes(self, props: dict[str, str]) -> list[str]:
        return self._parse_notes(props, "ноты сердца")

    def _parse_base_notes(self, props: dict[str, str]) -> list[str]:
        return self._parse_notes(props, "базовые ноты")

    def _get_shop_info(self, page: BeautifulSoup) -> PerfumeFromConcreteShop.ShopInfo:
        shop_info = PerfumeFromConcreteShop.ShopInfo(
            shop_name="Letu",
            domain="https://www.letu.ru",
            image_url=self._parse_image_url(page),
            variants=[],
        )
        current_item_variant = self._parse_current_item_variant(page)
        if current_item_variant:
            shop_info.variants.append(current_item_variant)

        other_variants_tags = page.find_all(
            "a", class_="product-detail-sku-volume sku-view-table__item-volume"
        )
        other_links: list[str] = []
        for other_variant_tag in other_variants_tags:
            link = other_variant_tag.get("href")
            if not isinstance(link, str):
                continue
            other_links.append(self._normalize_link("https://www.letu.ru", link))

        for link in other_links:
            other_page = get_page(link, use_playwright=True)
            if not other_page:
                continue
            other_variant = self._parse_current_item_variant(other_page)
            if other_variant:
                shop_info.variants.append(other_variant)
        return shop_info

    def _parse_image_url(self, page: BeautifulSoup) -> str:
        image_tag = page.find(
            "img", class_="product-detail-media-carousel__horizontal-item-img"
        )
        if not image_tag or not isinstance(image_tag, Tag):
            return ""
        src = image_tag.get("src")
        if not isinstance(src, str):
            return ""
        return self._normalize_link("https://www.letu.ru", src)

    def _normalize_link(self, domain: str, link: str) -> str:
        if link.startswith("/"):
            return domain + link
        return link

    def _parse_current_item_variant(
        self, page: BeautifulSoup
    ) -> PerfumeFromConcreteShop.ShopInfo.VolumeWithPrices | None:
        volume = self._parse_current_item_volume(page)
        if not volume:
            return None
        price = self._parse_current_item_price(page)
        if not price:
            return None
        link = self._parse_current_item_link(page)
        if not link:
            return None
        return PerfumeFromConcreteShop.ShopInfo.VolumeWithPrices(volume, price, link)

    def _parse_current_item_volume(self, page: BeautifulSoup) -> int | None:
        header = self._parse_header(page).strip().split(",")
        if not header:
            return None
        volume_text = header[-1].strip()
        volume_match = re.search(r"\d+", volume_text)
        if not volume_match:
            return None
        volume_value = volume_match.group(0)
        try:
            return int(volume_value)
        except ValueError:
            return None

    def _parse_current_item_price(self, page: BeautifulSoup) -> int | None:
        price_tag = page.find("span", class_="price-block__price-current")
        if not price_tag:
            price_tag = page.find(
                "p", class_="product-detail-price-action-block__price--current"
            )

        if not price_tag:
            return None
        price_text = price_tag.get_text(strip=True)
        price_digits = "".join(c for c in price_text if c.isdigit())
        if not price_digits:
            return None
        try:
            return int(price_digits)
        except ValueError:
            return None

    def _parse_current_item_link(self, page: BeautifulSoup) -> str:
        canonical_tag = page.find("link", rel="canonical")
        if canonical_tag and isinstance(canonical_tag, Tag):
            href = canonical_tag.get("href")
            if isinstance(href, str):
                return href

        og_url_tag = page.find("meta", property="og:url")
        if og_url_tag and isinstance(og_url_tag, Tag):
            content = og_url_tag.get("content")
            if isinstance(content, str):
                return content
        return ""
