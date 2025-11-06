from .perfume import PerfumeFromConcreteShop


class PerfumeWithUnitedShops:
    brand: str
    name: str
    sex: str
    properties: PerfumeFromConcreteShop.PerfumeProperties
    shops: list[PerfumeFromConcreteShop.ShopInfo]

    def __init__(self, perfume: PerfumeFromConcreteShop):
        self.brand = perfume.brand
        self.name = perfume.name
        self.sex = perfume.sex
        self.properties = perfume.properties
        self.shops = [perfume.shop_info]

    def to_dict(
        self,
    ) -> dict[
        str,
        str
        | dict[str, str | list[str]]
        | list[dict[str, str | list[dict[str, str | int]]]],
    ]:
        return {
            "brand": self.brand,
            "name": self.name,
            "sex": self.sex,
            "properties": self.properties.to_dict(),
            "shops": [shop_info.to_dict() for shop_info in self.shops],
        }
