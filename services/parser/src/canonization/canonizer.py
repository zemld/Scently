import csv
from pathlib import Path

from .levenshtein_distance import get_levenshtein_distance


class Canonizer:
    _levenshtein_diff = 2
    mapping: dict[str, str]

    def __init__(self, mapping_path: Path, levenshtein_diff: int = 2):
        self.mapping = self._load_mapping(mapping_path)
        self._levenshtein_diff = levenshtein_diff

    def _load_mapping(self, path: Path) -> dict[str, str]:
        mapping = {}
        for file in path.glob("*.csv"):
            with open(file) as f:
                rows = csv.reader(f, delimiter=";")
                for row in rows:
                    try:
                        key, value = row[0].strip(), row[1].strip()
                        mapping[key] = value
                    except Exception:
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
        lev_diff = self._levenshtein_diff + 1
        result = None
        for key in self.mapping.keys():
            lev_dist = get_levenshtein_distance(item, key)
            if lev_dist < lev_diff and lev_dist <= self._levenshtein_diff:
                result = self.mapping[key]
                lev_diff = lev_dist
        return result

    def canonize(self, word: str) -> str:
        exact = self._canonize_with_exact(word)
        if exact:
            return exact
        prefix = self._canonize_with_prefix(word)
        if prefix:
            return prefix
        levenshtein = self._canonize_with_levenshtein(word)
        if levenshtein:
            return levenshtein
        return ""
