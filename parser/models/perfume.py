class Perfume:
    brand: str
    name: str
    perfume_type: str
    sex: str
    family: str
    upper_notes: list[str]
    middle_notes: list[str]
    base_notes: list[str]
    volume: int
    link: str

    def __init__(
        self,
        brand: str = "",
        name: str = "",
        perfume_type: str = "",
        sex: str = "unisex",
        family: str = "",
        upper_notes: list[str] = [],
        middle_notes: list[str] = [],
        base_notes: list[str] = [],
        volume: int = 0,
        link: str = "",
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
        self.link = link

    def _repr_property(self, name: str, value: str) -> str:
        return f"{name}={value if value else 'Unknown'}"

    def __repr__(self):
        return (
            f"Perfume(\n\t{self._repr_property('brand', self.brand)},\n"
            f"\t{self._repr_property('name', self.name)},\n"
            f"\t{self._repr_property('perfume_type', self.perfume_type)},\n"
            f"\t{self._repr_property('sex', self.sex)},\n"
            f"\t{self._repr_property('family', self.family)},\n"
            f"\t{self._repr_property('upper_notes', self.upper_notes)},\n"
            f"\t{self._repr_property('middle_notes', self.middle_notes)},\n"
            f"\t{self._repr_property('base_notes', self.base_notes)},\n"
            f"\t{self._repr_property('volume', self.volume)},\n"
            f"\t{self._repr_property('link', self.link)}\n)\n"
        )

    def to_dict(self):
        return {
            "brand": self.brand,
            "name": self.name,
            "type": self.perfume_type,
            "sex": self.sex,
            "family": self.family,
            "upper_notes": self.upper_notes,
            "middle_notes": self.middle_notes,
            "base_notes": self.base_notes,
            "volume": self.volume,
            "link": self.link,
        }
