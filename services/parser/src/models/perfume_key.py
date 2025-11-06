from .perfume import PerfumeFromConcreteShop


class PerfumeKey:
    brand: str
    name: str
    sex: str

    def __init__(self, perfume: PerfumeFromConcreteShop):
        self.brand = perfume.brand
        self.name = perfume.name
        self.sex = perfume.sex

    def __hash__(self) -> int:
        return hash((self.brand, self.name, self.sex))

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, self.__class__):
            return False
        return (
            self.brand == other.brand
            and self.name == other.name
            and self.sex == other.sex
        )
