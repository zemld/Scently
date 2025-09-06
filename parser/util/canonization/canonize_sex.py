CANONIZATION_MAP = {"муж": "male", "жен": "female", "мальч": "male", "девоч": "female"}
UNISEX = "unisex"


def _canonize_with_exact(sex: str) -> str | None:
    return CANONIZATION_MAP.get(sex)


def _canonize_with_prefix(sex: str) -> str | None:
    for key in CANONIZATION_MAP.keys():
        if sex.startswith(key):
            return CANONIZATION_MAP[key]
    return None


def canonize_sex(sex: str) -> str:
    for word in sex.split():
        exact = _canonize_with_exact(word)
        if exact:
            return exact
        canonized = _canonize_with_prefix(word)
        if canonized:
            return canonized
    return UNISEX
