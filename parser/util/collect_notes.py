import re
import json


def read_file() -> str:
    with open("goldapple_perfumes.json", "r") as file:
        data = json.load(file)
    return data


def get_notes(data: str, note_field: str) -> set[str]:
    notes = set()
    for item in data:
        notes.update(item.get(note_field))
    return notes


def store_notes(notes: set[str], filename: str) -> None:
    with open("notes/" + filename, "w") as file:
        for note in notes:
            file.write(f"{note}\n")


if __name__ == "__main__":
    file = read_file()
    notes = set()
    notes.update(get_notes(file, "upper_notes"))
    notes.update(get_notes(file, "middle_notes"))
    notes.update(get_notes(file, "base_notes"))
    store_notes(notes, "raw_notes.csv")
