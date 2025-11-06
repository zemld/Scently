import argparse
import json
from collections.abc import Callable
from pathlib import Path

from src.canonization import Canonizer, NoteCanonizer
from src.scraping import Scrapper
from src.scraping.gold_apple import GoldApplePageParser, GoldAppleScrapper
from src.scraping.letu import LetuPageParser, LetuScrapper
from src.scraping.randewoo import RandewooPageParser, RandewooScrapper

PARSERS = {
    "goldapple": GoldApplePageParser(
        type_canonizer=Canonizer(Path.cwd() / "data/types"),
        sex_canonizer=Canonizer(Path.cwd() / "data/sex"),
        family_canonizer=Canonizer(Path.cwd() / "data/families"),
        notes_canonizer=NoteCanonizer(Path.cwd() / "data/notes"),
    ),
    "randewoo": RandewooPageParser(
        type_canonizer=Canonizer(Path.cwd() / "data/types"),
        sex_canonizer=Canonizer(Path.cwd() / "data/sex"),
        family_canonizer=Canonizer(Path.cwd() / "data/families"),
        notes_canonizer=NoteCanonizer(Path.cwd() / "data/notes"),
    ),
    "letu": LetuPageParser(
        type_canonizer=Canonizer(Path.cwd() / "data/types"),
        sex_canonizer=Canonizer(Path.cwd() / "data/sex"),
        family_canonizer=Canonizer(Path.cwd() / "data/families"),
        notes_canonizer=NoteCanonizer(Path.cwd() / "data/notes"),
    ),
}

SCRAPPER_FACTORIES: dict[str, Callable[[], Scrapper]] = {
    "goldapple": lambda: GoldAppleScrapper(PARSERS["goldapple"], max_pages=2),
    "randewoo": lambda: RandewooScrapper("https://randewoo.ru", PARSERS["randewoo"]),
    "letu": lambda: LetuScrapper("https://www.letu.ru", PARSERS["letu"]),
}

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--shop", required=True, help="Название магазина")
    args = parser.parse_args()
    shop_name = args.shop

    scrapper_factory = SCRAPPER_FACTORIES.get(shop_name)
    if not scrapper_factory:
        available_shops = ", ".join(sorted(SCRAPPER_FACTORIES))
        raise ValueError(
            f"Неизвестный магазин '{shop_name}'. Доступные магазины: {available_shops}"
        )

    scrapper = scrapper_factory()

    perfumes = scrapper.scrap_page(0)
    with open(f"data/collected_perfumes/{shop_name}_perfumes.json", "w") as f:
        json.dump(
            [perfume.to_dict() for perfume in perfumes],
            f,
            indent=4,
            ensure_ascii=False,
        )
