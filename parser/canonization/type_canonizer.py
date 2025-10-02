from canonization.canonizer import Canonizer
from util.levenshtein_distance import get_levenshtein_distance


class TypeCanonizer(Canonizer):
    def canonize(self, perfume_type: str) -> str | None:
        exact = super()._canonize_with_exact(perfume_type)
        if exact:
            return exact
        lev_diff = super()._levenshtein_diff + 1
        result = None
        for key in self.mapping.keys():
            lev_dist = get_levenshtein_distance(perfume_type, key)
            if lev_dist < lev_diff:
                lev_diff = lev_dist
                result = self.mapping[key]
        if not result:
            raise ValueError(f"Cannot canonize type: {perfume_type}")
        return result
