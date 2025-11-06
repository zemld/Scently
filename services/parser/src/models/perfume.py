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
