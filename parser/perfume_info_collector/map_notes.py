import csv
from pathlib import Path

NOTES_MAP = Path("data/notes.csv")


def _parse_row(row: str) -> list[str] | None:
    stripped = row.strip("'[]")
    return stripped.split(";")


def load_notes_map() -> dict[str, str]:
    notes_map = {}
    with open(NOTES_MAP, "r", encoding="utf-8") as f:
        reader = csv.reader(f)
        for row in reader:
            parsed_row = _parse_row(row[0])
            notes_map[parsed_row[0]] = parsed_row[1]
    return notes_map


if __name__ == "__main__":
    notes_map = load_notes_map()
    print(notes_map)
