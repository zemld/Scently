from typing import cast


class PerfumeFromConcreteShop:
    class PerfumeProperties:
        perfume_type: str
        family: list[str]
        upper_notes: list[str]
        middle_notes: list[str]
        base_notes: list[str]

        def __init__(
            self,
            perfume_type: str,
            family: list[str],
            upper_notes: list[str],
            middle_notes: list[str],
            base_notes: list[str],
        ):
            self.perfume_type = perfume_type
            self.family = family
            self.upper_notes = upper_notes
            self.middle_notes = middle_notes
            self.base_notes = base_notes

        def to_dict(self) -> dict[str, str | list[str]]:
            return {
                "perfume_type": self.perfume_type,
                "family": self.family,
                "upper_notes": self.upper_notes,
                "middle_notes": self.middle_notes,
                "base_notes": self.base_notes,
            }

        @classmethod
        def from_dict(
            cls, data: dict[str, str | list[str]]
        ) -> "PerfumeFromConcreteShop.PerfumeProperties":
            perfume_type = data["perfume_type"]
            family = data["family"]
            upper_notes = data["upper_notes"]
            middle_notes = data["middle_notes"]
            base_notes = data["base_notes"]

            assert isinstance(perfume_type, str), "perfume_type must be str"
            assert isinstance(family, list), "family must be list[str]"
            assert isinstance(upper_notes, list), "upper_notes must be list[str]"
            assert isinstance(middle_notes, list), "middle_notes must be list[str]"
            assert isinstance(base_notes, list), "base_notes must be list[str]"

            return PerfumeFromConcreteShop.PerfumeProperties(
                perfume_type=perfume_type,
                family=family,
                upper_notes=upper_notes,
                middle_notes=middle_notes,
                base_notes=base_notes,
            )

    class ShopInfo:
        class VolumeWithPrices:
            volume: int
            cost: int
            link: str

            def __init__(self, volume: int, cost: int, link: str):
                self.volume = volume
                self.cost = cost
                self.link = link

            def to_dict(self) -> dict[str, int | str]:
                return {
                    "volume": self.volume,
                    "price": self.cost,
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
        shop_link: str
        image_url: str
        volumes_with_prices: list[VolumeWithPrices]

        def __init__(
            self,
            shop_name: str,
            shop_link: str,
            image_url: str,
            volumes_with_prices: list[VolumeWithPrices],
        ):
            self.shop_name = shop_name
            self.shop_link = shop_link
            self.image_url = image_url
            self.volumes_with_prices = volumes_with_prices

        def to_dict(self) -> dict[str, str | list[dict[str, str | int]]]:
            return {
                "shop_name": self.shop_name,
                "shop_link": self.shop_link,
                "image_url": self.image_url,
                "volumes_with_prices": [v.to_dict() for v in self.volumes_with_prices],
            }

        @classmethod
        def from_dict(
            cls, data: dict[str, str | list[dict[str, str | int]]]
        ) -> "PerfumeFromConcreteShop.ShopInfo":
            shop_name = data["shop_name"]
            shop_link = data["shop_link"]
            image_url = data["image_url"]
            volumes_with_prices = data["volumes_with_prices"]

            assert isinstance(shop_name, str), "shop_name must be str"
            assert isinstance(shop_link, str), "shop_link must be str"
            assert isinstance(image_url, str), "image_url must be str"
            assert isinstance(
                volumes_with_prices, list
            ), "volumes_with_prices must be list"

            return PerfumeFromConcreteShop.ShopInfo(
                shop_name=shop_name,
                shop_link=shop_link,
                image_url=image_url,
                volumes_with_prices=[
                    PerfumeFromConcreteShop.ShopInfo.VolumeWithPrices.from_dict(v)
                    for v in volumes_with_prices
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
            middle_notes=[],
            base_notes=[],
        )
        self.shop_info = shop_info or self.ShopInfo(
            shop_name="",
            shop_link="",
            image_url="",
            volumes_with_prices=[],
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
