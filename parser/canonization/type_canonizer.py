from canonization.canonizer import Canonizer
from canonization.levenshtein_distance import get_levenshtein_distance


class TypeCanonizer(Canonizer):
    def canonize(self, perfume_type: list[str]) -> str | None:
        if not perfume_type:
            return None
        perfume_type_word = perfume_type[0]
        exact = super()._canonize_with_exact(perfume_type_word)
        if exact:
            return exact
        lev_diff = self._levenshtein_diff + 1
        result = None
        for key in self.mapping.keys():
            lev_dist = get_levenshtein_distance(perfume_type_word, key)
            if lev_dist < lev_diff:
                lev_diff = lev_dist
                result = self.mapping[key]
        return result
