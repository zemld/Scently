from util.canonization.map_notes import load_notes_map
from util.levenshtein_distance import get_levenshtein_distance
from pathlib import Path

NOTES_MAP = load_notes_map(Path("data/notes"))

LEVENSHTEIN_THRESHOLD = 4
DIFF = 2

MEANINGLESS_WORDS = (
    "масло",
    "масла",
    "ноты",
    "нот",
    "аккорд",
    "аккорды",
    "цветок",
    "цветки",
    "цветов",
    "цвет",
    "цветы",
    "экстракт",
    "экстракты",
    "экстракта",
    "абсолют",
    "абсолюта",
    "эссенция",
    "эссенции",
    "породы",
    "порода",
    "пород",
)
ADJECTIVE_ENDINGS = (
    "ий",
    "ый",
    "ая",
    "ое",
    "ие",
    "ые",
    "ого",
    "ой",
    "ых",
    "ому",
    "им",
    "ым",
    "ую",
    "ым",
    "ими",
    "ыми",
    "ом",
    "их",
    "ых",
)


def _canonize_with_exact(note: str) -> str | None:
    return NOTES_MAP.get(note)


def _canonize_with_prefix(note: str) -> str | None:
    diff = 1000
    result = None
    for key in NOTES_MAP.keys():
        canonized = NOTES_MAP[key] if note.startswith(key) else None
        if canonized and len(note) - len(key) < diff:
            result = canonized
            diff = len(note) - len(key)
    return result


def _canonize_with_levenshtein(note: str) -> str | None:
    lev_diff = DIFF + 1
    result = None
    if len(note) <= LEVENSHTEIN_THRESHOLD:
        return None
    for key in NOTES_MAP.keys():
        lev_dist = get_levenshtein_distance(note, key)
        if lev_dist < lev_diff:
            result = NOTES_MAP[key]
            lev_diff = lev_dist
    return result


def _canonize_with_essential_word(word: str) -> str | None:
    exact = _canonize_with_exact(word)
    if exact:
        return exact
    canonized = _canonize_with_prefix(word)
    if canonized:
        return canonized
    return _canonize_with_levenshtein(word)


def _canonize_by_words(note: list[str]) -> str | None:
    for word in note:
        canonized = _canonize_with_essential_word(word)
        if canonized:
            return canonized
    return None


def _preprocess_note(note_words: list[str]) -> list[str]:
    return [word.strip(",. ") for word in note_words if word.strip(",. ")]


def _divide_with_adjectives(words: list[str]) -> tuple[list[str], list[str]]:
    without_adjectives = [
        word for word in words if not word.endswith(ADJECTIVE_ENDINGS)
    ]
    adjectives = [word for word in words if word.endswith(ADJECTIVE_ENDINGS)]
    return without_adjectives, adjectives


def canonize_note(note: str) -> str | None:
    exact = NOTES_MAP.get(note)
    if exact:
        return exact
    words = note.split()
    essential_words = [
        word for word in _preprocess_note(words) if word not in MEANINGLESS_WORDS
    ]
    without_adjectives, adjectives = _divide_with_adjectives(essential_words)
    canonized = _canonize_by_words(without_adjectives)
    if canonized:
        return canonized
    canonized = _canonize_by_words(essential_words)
    if canonized:
        return canonized
    return _canonize_with_levenshtein(note)
