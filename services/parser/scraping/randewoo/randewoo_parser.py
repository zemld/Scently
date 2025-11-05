import re

from bs4 import BeautifulSoup, Tag

from models import Perfume

from ..page_parser import PageParser


def _script_contains_product_variable(text: str | None) -> bool:
    return isinstance(text, str) and "var product" in text


class RandewooPageParser(PageParser):
    def _parse_brand(self, page: BeautifulSoup) -> str:
        brand_tags = page.find_all("div", class_="b-header__mainTitle")
        if not brand_tags:
            return ""
        return str(brand_tags[0].get_text(strip=True))

    def _parse_name(self, page: BeautifulSoup) -> str:
        name_tags = page.find_all("div", class_="b-header__subtitle")
        if not name_tags:
            return ""
        return str(name_tags[0].get_text(strip=True))

    def _parse_type(self, page: BeautifulSoup) -> str:
        type_tags = page.find_all("a", class_="s-productType__titleText")
        if not type_tags:
            return self._fallback_parse_type(page)
        type_text = str(type_tags[0].get_text(separator=" ", strip=True))
        if not type_text:
            return self._fallback_parse_type(page)
        normalized_type = self._normalize_type_text(type_text)
        if normalized_type:
            return normalized_type
        return self._fallback_parse_type(page)

    def _fallback_parse_type(self, page: BeautifulSoup) -> str:
        short_description = page.find(
            "h3", class_="p-productCard__shortDescriptionText"
        )
        if short_description:
            type_text = str(short_description.get_text(separator=" ", strip=True))
            normalized_type = self._normalize_type_text(type_text)
            if normalized_type:
                return normalized_type

        offer_meta = page.select_one('[itemprop="offers"] meta[itemprop="name"]')
        if isinstance(offer_meta, Tag):
            offer_content = offer_meta.get("content")
            if isinstance(offer_content, str):
                normalized_type = self._normalize_type_text(offer_content)
                if normalized_type:
                    return normalized_type

        for script_tag in page.find_all("script"):
            if not _script_contains_product_variable(script_tag.string):
                continue
            script_content = script_tag.string
            if not isinstance(script_content, str):
                continue
            subtitle_match = re.search(r'"subtitle"\s*:\s*"([^"]+)"', script_content)
            if subtitle_match:
                normalized_type = self._normalize_type_text(subtitle_match.group(1))
                if normalized_type:
                    return normalized_type

        return ""

    def _normalize_type_text(self, raw_text: str) -> str:
        prepared_text = raw_text.replace("\xa0", " ").strip().lower()
        if not prepared_text:
            return ""

        volume_match = re.search(r"(\d+)\s*мл", prepared_text)
        if volume_match:
            prepared_text = prepared_text[: volume_match.start()].strip()

        if not prepared_text:
            return ""

        if self._type_canonizer:
            canonized = self._type_canonizer.canonize(prepared_text)
            if canonized:
                return canonized

        return prepared_text

    def _extract_volume_from_text(self, raw_text: str) -> int:
        normalized_text = raw_text.replace("\xa0", " ").lower()
        volume_match = re.search(r"(\d+)\s*мл", normalized_text)
        if not volume_match:
            return 0
        try:
            return int(volume_match.group(1))
        except ValueError:
            return 0

    def _extract_volume(self, page: BeautifulSoup) -> int:
        short_description = page.find(
            "h3", class_="p-productCard__shortDescriptionText"
        )
        if not short_description:
            return 0
        volume = self._extract_volume_from_text(
            str(short_description.get_text(separator=" ", strip=True))
        )
        if volume:
            return volume

        offer_meta = page.select_one('[itemprop="offers"] meta[itemprop="name"]')
        if isinstance(offer_meta, Tag):
            offer_content = offer_meta.get("content")
            if isinstance(offer_content, str):
                volume = self._extract_volume_from_text(offer_content)
                if volume:
                    return volume

        for script_tag in page.find_all("script"):
            if not _script_contains_product_variable(script_tag.string):
                continue
            script_content = script_tag.string
            if not isinstance(script_content, str):
                continue
            subtitle_match = re.search(r'"subtitle"\s*:\s*"([^"]+)"', script_content)
            if subtitle_match:
                volume = self._extract_volume_from_text(subtitle_match.group(1))
                if volume:
                    return volume

        return 0

    def _extract_price(self, page: BeautifulSoup) -> int | None:
        price_tag = page.find("strong", class_="b-productSummary__priceNew")
        if price_tag:
            raw_price = str(price_tag.get_text(strip=True))
            price_digits = re.sub(r"\D", "", raw_price)
            if price_digits:
                try:
                    return int(price_digits)
                except ValueError:
                    return None

        price_meta = page.select_one('[itemprop="offers"] meta[itemprop="price"]')
        if isinstance(price_meta, Tag):
            price_content = price_meta.get("content")
            if isinstance(price_content, str):
                price_digits = re.sub(r"\D", "", price_content)
                if price_digits:
                    try:
                        return int(price_digits)
                    except ValueError:
                        return None

        for script_tag in page.find_all("script"):
            if not _script_contains_product_variable(script_tag.string):
                continue
            script_content = script_tag.string
            if not isinstance(script_content, str):
                continue
            price_match = re.search(r'"price"\s*:\s*(\d+)', script_content)
            if price_match:
                try:
                    return int(price_match.group(1))
                except ValueError:
                    return None

        return None

    def _parse_sex(self, props: dict[str, str]) -> str:
        raw_sex = props.get("пол")
        if not raw_sex:
            return ""
        if " " in raw_sex:
            raw_sex = raw_sex.split()[-1]
        if self._sex_canonizer:
            result = self._sex_canonizer.canonize(raw_sex)
            return result if result else raw_sex
        return raw_sex

    def _parse_families(self, props: dict[str, str]) -> list[str]:
        families = props.get("группы")
        if not families:
            return []
        if not self._family_canonizer:
            return [family.strip() for family in families.split(",") if family.strip()]
        canonized_families = [
            self._family_canonizer.canonize(family.strip())
            for family in families.split(",")
        ]
        not_none_families = [family for family in canonized_families if family]
        return list(set(not_none_families))

    def _parse_notes(self, props: dict[str, str], key: str) -> list[str]:
        notes = props.get(key)
        if not notes:
            return []
        if not self._notes_canonizer:
            return [note.strip() for note in notes.split(",") if note.strip()]
        canonized_notes = [
            self._notes_canonizer.canonize(note.strip()) for note in notes.split(",")
        ]
        return [note for note in canonized_notes if note]

    def _parse_upper_notes(self, props: dict[str, str]) -> list[str]:
        return self._parse_notes(props, "верхние ноты")

    def _parse_middle_notes(self, props: dict[str, str]) -> list[str]:
        return self._parse_notes(props, "средние ноты")

    def _parse_base_notes(self, props: dict[str, str]) -> list[str]:
        return self._parse_notes(props, "базовые ноты")

    def _parse_image_url(self, page: BeautifulSoup) -> str:
        image_tag = page.find(
            "img", class_="js-main-product-image s-productItem__imgMain"
        )
        if not image_tag:
            return ""
        if not isinstance(image_tag, Tag):
            return ""
        src = image_tag.get("src")
        return src if isinstance(src, str) else ""

    def _parse_props(self, page: BeautifulSoup) -> dict[str, str]:
        all_prop_tags = page.find_all("dl", class_="dl")
        if not all_prop_tags:
            return {}

        current_prop_tags = all_prop_tags[0].find_all("div")
        if not current_prop_tags:
            return {}

        props: dict[str, str] = {}
        for prop_tag in current_prop_tags:
            key_tag = prop_tag.find("dt")
            value_tag = prop_tag.find("dd")
            if not key_tag or not value_tag:
                continue

            key = key_tag.get_text(strip=True).lower()
            if not key:
                continue

            link_texts = [
                anchor.get_text(strip=True)
                for anchor in value_tag.find_all("a")
                if anchor.get_text(strip=True)
            ]
            if link_texts:
                value = ", ".join(link_texts)
            else:
                value = value_tag.get_text(separator=" ", strip=True)

            value = value.strip()
            if not value:
                continue

            props[key] = value.lower()

        return props

    def _get_shop_info(self, page: BeautifulSoup) -> Perfume.ShopInfo:
        shop_info = Perfume.ShopInfo(
            shop_name="Randewoo",
            shop_link="https://randewoo.ru",
            volumes_with_prices=[],
        )
        volumes_with_prices_tags = page.find_all("div", class_="s-productType__main")
        if not volumes_with_prices_tags:
            price_value = self._extract_price(page)
            if price_value is None:
                return shop_info
            volume_value = self._extract_volume(page)
            shop_info.volumes_with_prices.append(
                Perfume.ShopInfo.VolumeWithPrices(
                    volume_value,
                    price_value,
                    "",
                )
            )
            return shop_info
        for volumes_with_prices_tag in volumes_with_prices_tags:
            volume_anchor = volumes_with_prices_tag.find(
                "a",
                class_=(
                    "s-productType__titleText s-link "
                    "s-link--unbordered js-product-show_image"
                ),
            )
            if volume_anchor is None:
                continue
            raw_volume = volume_anchor.get_text(strip=True)
            volume_text = raw_volume.split()[-1].strip("мл ")
            try:
                volume_value = int(volume_text)
            except ValueError:
                continue
            cost_tag = volumes_with_prices_tag.find(
                "span", class_="s-productType__priceNewValue"
            )
            if cost_tag is None:
                continue
            raw_cost = str(cost_tag.get_text(strip=True))
            price_digits = re.sub(r"\D", "", raw_cost)
            if not price_digits:
                continue
            shop_info.volumes_with_prices.append(
                Perfume.ShopInfo.VolumeWithPrices(
                    int(volume_text), int(price_digits), ""
                )
            )

        return shop_info
