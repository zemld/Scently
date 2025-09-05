import csv
from pathlib import Path

NOTES_MAP = Path("data/notes")


def _parse_row(row: str) -> list[str] | None:
    stripped = row.strip("'[]")
    return stripped.split(";")


def load_notes_map(path: Path) -> dict[str, str]:
    notes_map = {}
    for file in path.glob("*.csv"):
        with open(file, "r") as f:
            reader = csv.reader(f, delimiter=";")
            for raw in reader:
                if not raw or len(raw) < 2:
                    continue
                notes_map[raw[0]] = raw[1]
    return notes_map


if __name__ == "__main__":
    m = load_notes_map(NOTES_MAP)
    print(len(m))
