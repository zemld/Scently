import json
from pathlib import Path

if __name__ == "__main__":
    perfumes = []
    for file in Path("data/collected_perfumes").glob("*.json"):
        with open(file) as f:
            perfumes.extend(json.load(f))
    print(len(perfumes))
    with open("data/all_perfumes.json", "w") as f:
        json.dump(perfumes, f)
