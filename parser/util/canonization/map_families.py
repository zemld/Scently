import csv
from pathlib import Path

FAMILIES_PATH = Path("data/families/families.csv")


def load_families_map(path: Path) -> dict[str, str]:
    families_map = {}
    with open(path, "r", encoding="utf-8") as f:
        reader = csv.reader(f, delimiter=";")
        for row in reader:
            families_map[row[0].strip().lower()] = row[1].strip().lower()
    return families_map
