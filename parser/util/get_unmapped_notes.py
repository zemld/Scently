from util.map_notes import load_notes_map
from util.canonize_note import canonize_note
from pathlib import Path
import json

MAP_PATH = Path("data/notes")


def load_raw_notes():
    raw_notes = set()
    with open("goldapple_perfumes.json", "r", encoding="utf-8") as file:
        data = json.load(file)
        for item in data:
            raw_notes.update(item["upper_notes"])
            raw_notes.update(item["middle_notes"])
            raw_notes.update(item["base_notes"])
    return raw_notes


if __name__ == "__main__":
    raw_notes = load_raw_notes()
    notes_map = load_notes_map(MAP_PATH)
    mapped = {}
    unmapped = []
    for note in raw_notes:
        canonization = canonize_note(note)
        if canonization:
            mapped[note] = canonization
        else:
            unmapped.append(note)
    with open("unmapped_notes.txt", "w", encoding="utf-8") as file:
        for item in unmapped:
            file.write(f"{item}\n")
    with open("mapped_notes.json", "w", encoding="utf-8") as file:
        json.dump(mapped, file, ensure_ascii=False, indent=4)
    print(f"Total raw notes: {len(raw_notes)}")
    print(f"Mapped notes: {len(mapped)}")
    print(f"Unmapped notes: {len(unmapped)}")
    print(f"Mapping ratio: {len(mapped) / len(raw_notes):.2%}")
