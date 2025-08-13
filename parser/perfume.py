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
