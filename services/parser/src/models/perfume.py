from typing import cast


class PerfumeFromConcreteShop:
    class PerfumeProperties:
        perfume_type: str
        family: list[str]
        upper_notes: list[str]
        core_notes: list[str]
        base_notes: list[str]

        def __init__(
            self,
            perfume_type: str,
            family: list[str],
            upper_notes: list[str],
            core_notes: list[str],
            base_notes: list[str],
        ):
            self.perfume_type = perfume_type
            self.family = family
            self.upper_notes = upper_notes
            self.core_notes = core_notes
            self.base_notes = base_notes

        def to_dict(self) -> dict[str, str | list[str]]:
            return {
                "perfume_type": self.perfume_type,
                "family": self.family,
                "upper_notes": self.upper_notes,
                "core_notes": self.core_notes,
                "base_notes": self.base_notes,
            }

        @classmethod
        def from_dict(
            cls, data: dict[str, str | list[str]]
        ) -> "PerfumeFromConcreteShop.PerfumeProperties":
            perfume_type = data["perfume_type"]
            family = data["family"]
            upper_notes = data["upper_notes"]
            core_notes = data["core_notes"]
            base_notes = data["base_notes"]

            assert isinstance(perfume_type, str), "perfume_type must be str"
            assert isinstance(family, list), "family must be list[str]"
            assert isinstance(upper_notes, list), "upper_notes must be list[str]"
            assert isinstance(core_notes, list), "core_notes must be list[str]"
            assert isinstance(base_notes, list), "base_notes must be list[str]"

            return PerfumeFromConcreteShop.PerfumeProperties(
                perfume_type=perfume_type,
                family=family,
                upper_notes=upper_notes,
                core_notes=core_notes,
                base_notes=base_notes,
            )

    class ShopInfo:
        class VolumeWithPrices:
            volume: int
            price: int
            link: str

            def __init__(self, volume: int, cost: int, link: str):
                self.volume = volume
                self.price = cost
                self.link = link

            def to_dict(self) -> dict[str, int | str]:
                return {
                    "volume": self.volume,
                    "price": self.price,
                    "link": self.link,
                }

            @classmethod
            def from_dict(
                cls, data: dict[str, int | str]
            ) -> "PerfumeFromConcreteShop.ShopInfo.VolumeWithPrices":
                volume = data["volume"]
                cost = data["price"]
                link = data["link"]

                assert isinstance(volume, int), "volume must be int"
                assert isinstance(cost, int), "cost must be int"
                assert isinstance(link, str), "link must be str"

                return PerfumeFromConcreteShop.ShopInfo.VolumeWithPrices(
                    volume=volume,
                    cost=cost,
                    link=link,
                )

        shop_name: str
        domain: str
        image_url: str
        variants: list[VolumeWithPrices]

        def __init__(
            self,
            shop_name: str,
            domain: str,
            image_url: str,
            variants: list[VolumeWithPrices],
        ):
            self.shop_name = shop_name
            self.domain = domain
            self.image_url = image_url
            self.variants = variants

        def to_dict(self) -> dict[str, str | list[dict[str, str | int]]]:
            return {
                "shop_name": self.shop_name,
                "domain": self.domain,
                "image_url": self.image_url,
                "variants": [v.to_dict() for v in self.variants],
            }

        @classmethod
        def from_dict(
            cls, data: dict[str, str | list[dict[str, str | int]]]
        ) -> "PerfumeFromConcreteShop.ShopInfo":
            shop_name = data["shop_name"]
            domain = data["domain"]
            image_url = data["image_url"]
            variants = data["variants"]

            assert isinstance(shop_name, str), "shop_name must be str"
            assert isinstance(domain, str), "domain must be str"
            assert isinstance(image_url, str), "image_url must be str"
            assert isinstance(variants, list), "volumes_with_prices must be list"

            return PerfumeFromConcreteShop.ShopInfo(
                shop_name=shop_name,
                domain=domain,
                image_url=image_url,
                variants=[
                    PerfumeFromConcreteShop.ShopInfo.VolumeWithPrices.from_dict(v)
                    for v in variants
                ],
            )

    brand: str
    name: str
    sex: str
    properties: PerfumeProperties
    shop_info: ShopInfo

    def __init__(
        self,
        brand: str = "",
        name: str = "",
        sex: str = "unisex",
        properties: PerfumeProperties | None = None,
        shop_info: ShopInfo | None = None,
    ):
        self.brand = brand
        self.name = name
        self.sex = sex
        self.properties = properties or self.PerfumeProperties(
            perfume_type="",
            family=[],
            upper_notes=[],
            core_notes=[],
            base_notes=[],
        )
        self.shop_info = shop_info or self.ShopInfo(
            shop_name="",
            domain="",
            image_url="",
            variants=[],
        )

    def to_dict(
        self,
    ) -> dict[
        str,
        str | dict[str, str | list[str]] | dict[str, str | list[dict[str, str | int]]],
    ]:
        return {
            "brand": self.brand,
            "name": self.name,
            "sex": self.sex,
            "properties": self.properties.to_dict(),
            "shop_info": self.shop_info.to_dict(),
        }

    @classmethod
    def from_dict(
        cls,
        data: dict[
            str,
            str
            | dict[str, str | list[str]]
            | dict[str, str | list[dict[str, str | int]]],
        ],
    ) -> "PerfumeFromConcreteShop":
        brand = data["brand"]
        name = data["name"]
        sex = data["sex"]
        properties = data["properties"]
        shop_info = data["shop_info"]

        assert isinstance(brand, str), "brand must be str"
        assert isinstance(name, str), "name must be str"
        assert isinstance(sex, str), "sex must be str"
        assert isinstance(properties, dict), "properties must be dict"
        assert isinstance(shop_info, dict), "shop_info must be dict"

        properties_dict = cast(dict[str, str | list[str]], properties)
        shop_info_dict = cast(dict[str, str | list[dict[str, str | int]]], shop_info)

        return PerfumeFromConcreteShop(
            brand=brand,
            name=name,
            sex=sex,
            properties=PerfumeFromConcreteShop.PerfumeProperties.from_dict(
                properties_dict
            ),
            shop_info=PerfumeFromConcreteShop.ShopInfo.from_dict(shop_info_dict),
        )
