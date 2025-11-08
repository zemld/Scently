import json
import shutil
from pathlib import Path

from src.app import get_all_perfumes, unite_perfumes
from src.canonization import Canonizer, NoteCanonizer
from src.scraping.gold_apple import GoldApplePageParser, GoldAppleScrapper
from src.scraping.letu import LetuPageParser, LetuScrapper
from src.scraping.randewoo import RandewooPageParser, RandewooScrapper
from src.scraping.scrapper import Scrapper

UPLOAD_PERFUME_INFO = "http://perfume:8089/v1/perfumes/update"

TYPE_CANONIZER = Canonizer(Path.cwd() / "data/types")
SEX_CANONIZER = Canonizer(Path.cwd() / "data/sex")
FAMILY_CANONIZER = Canonizer(Path.cwd() / "data/families")
NOTES_CANONIZER = NoteCanonizer(Path.cwd() / "data/notes")


SCRAPPERS = {
    "goldapple": GoldAppleScrapper(
        page_parser=GoldApplePageParser(
            type_canonizer=TYPE_CANONIZER,
            sex_canonizer=SEX_CANONIZER,
            family_canonizer=FAMILY_CANONIZER,
            notes_canonizer=NOTES_CANONIZER,
        ),
    ),
    "randewoo": RandewooScrapper(
        "https://randewoo.ru",
        page_parser=RandewooPageParser(
            type_canonizer=TYPE_CANONIZER,
            sex_canonizer=SEX_CANONIZER,
            family_canonizer=FAMILY_CANONIZER,
            notes_canonizer=NOTES_CANONIZER,
        ),
    ),
    "letu": LetuScrapper(
        "https://www.letu.ru",
        page_parser=LetuPageParser(
            type_canonizer=TYPE_CANONIZER,
            sex_canonizer=SEX_CANONIZER,
            family_canonizer=FAMILY_CANONIZER,
            notes_canonizer=NOTES_CANONIZER,
        ),
    ),
}


def collect_and_store_perfumes(shop_name: str, scrapper: Scrapper) -> None:
    perfumes = scrapper.scrap_all_accuratly()
    with open(f"data/collected_perfumes/{shop_name}_perfumes.json", "w") as f:
        json.dump(
            [perfume.to_dict() for perfume in perfumes],
            f,
            indent=4,
            ensure_ascii=False,
        )


def collect_and_store_all_perfumes() -> None:
    for shop_name, scrapper in SCRAPPERS.items():
        collect_and_store_perfumes(shop_name, scrapper)
        print(f"Collected and stored perfumes for {shop_name}")


if __name__ == "__main__":
    collect_and_store_all_perfumes()
    perfumes = unite_perfumes(get_all_perfumes(Path.cwd() / "data/collected_perfumes"))
    with open(Path.cwd() / "data/collected_perfumes/all_perfumes.json", "w") as f:
        json.dump(
            [perfume.to_dict() for perfume in perfumes],
            f,
            indent=4,
            ensure_ascii=False,
        )
        backup_dir = Path.cwd() / "data/backups"
        if backup_dir.exists() and backup_dir.is_dir():
            shutil.rmtree(backup_dir)
