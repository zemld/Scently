from pydantic import BaseModel


class Perfume(BaseModel):
    brand: str
    name: str


class PerfumeOut(BaseModel):
    perfumes: list[Perfume]
