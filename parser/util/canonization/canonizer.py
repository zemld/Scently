from pathlib import Path
import csv
from util.levenshtein_distance import get_levenshtein_distance


class Canonizer:
    levenshtein_diff = 2
    mapping: dict[str, str]

    def __init__(self, mapping_path: Path):
        self.mapping = self._load_mapping(mapping_path)

    def _load_mapping(self, path: Path) -> dict[str, str]:
        mapping = {}
        for file in path.glob("*.csv"):
            with open(file, "r") as f:
                rows = csv.reader(f, delimiter=";")
                for row in rows:
                    try:
                        key, value = row[0].strip(), row[1].strip()
                        mapping[key] = value
                    except Exception as e:
                        continue
        return mapping

    def _canonize_with_exact(self, item: str) -> str | None:
        return self.mapping.get(item)

    def _canonize_with_prefix(self, item: str) -> str | None:
        diff = 1000
        result = None
        for key in self.mapping.keys():
            canonized = self.mapping[key] if item.startswith(key) else None
            if canonized and len(item) - len(key) < diff:
                result = canonized
                diff = len(item) - len(key)
        return result

    def _canonize_with_levenshtein(self, item: str) -> str | None:
        lev_diff = self.levenshtein_diff + 1
        result = None
        for key in self.mapping.keys():
            lev_dist = get_levenshtein_distance(item, key)
            if lev_dist < lev_diff and lev_dist <= self.levenshtein_diff:
                result = self.mapping[key]
                lev_diff = lev_dist
        return result

    def _canonize_with_word(self, word: str) -> str | None:
        exact = self._canonize_with_exact(word)
        if exact:
            return exact
        prefix = self._canonize_with_prefix(word)
        if prefix:
            return prefix
        return self._canonize_with_levenshtein(word)

    def canonize(self, words: list[str]) -> str | None:
        for word in words:
            canonized = self._canonize_with_word(word)
            if canonized:
                return canonized
        return None
