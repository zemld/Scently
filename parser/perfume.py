from enum import Enum


class Sex(Enum):
    FEMALE = "female"
    MALE = "male"
    UNISEX = "unisex"


class Perfume:
    brand: str
    name: str
    perfume_type: str
    sex: Sex
    family: str
    upper_notes: list[str]
    middle_notes: list[str]
    base_notes: list[str]
    volume: list[int]

    def __init__(
        self,
        brand: str = "",
        name: str = "",
        perfume_type: str = "",
        sex: Sex = Sex.UNISEX,
        family: str = "",
        upper_notes: list[str] = [],
        middle_notes: list[str] = [],
        base_notes: list[str] = [],
        volume: list[int] = [],
    ):
        self.brand = brand
        self.name = name
        self.perfume_type = perfume_type
        self.sex = sex
        self.family = family
        self.upper_notes = upper_notes
        self.middle_notes = middle_notes
        self.base_notes = base_notes
        self.volume = volume

    def _repr_property(self, name: str, value: str) -> str:
        return f"{name}={value if value else 'Unknown'}"

    def __repr__(self):
        return (
            f"Perfume({self._repr_property('brand', self.brand)}, "
            f"{self._repr_property('name', self.name)}, "
            f"{self._repr_property('perfume_type', self.perfume_type)}, "
            f"{self._repr_property('sex', self.sex)}, "
            f"{self._repr_property('family', self.family)}, "
            f"{self._repr_property('upper_notes', self.upper_notes)}, "
            f"{self._repr_property('middle_notes', self.middle_notes)}, "
            f"{self._repr_property('base_notes', self.base_notes)}, "
            f"{self._repr_property('volume', self.volume)})"
        )
