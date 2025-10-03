from canonization.canonizer import Canonizer


class NoteCanonizer(Canonizer):
    _meaningless_words = (
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
    _adjective_endings = (
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

    def _preprocess_note(self, note_words: list[str]) -> list[str]:
        return [word.strip(",. \n") for word in note_words if word.strip(",. ")]

    def _divide_with_adjectives(self, words: list[str]) -> tuple[list[str], list[str]]:
        without_adjectives = [
            word for word in words if not word.endswith(self._adjective_endings)
        ]
        adjectives = [word for word in words if word.endswith(self._adjective_endings)]
        return without_adjectives, adjectives

    def canonize(self, notes: list[str]) -> list[str]:
        canonized = []
        for note in notes:
            canonized.append(self._canonize_note(note))
        return canonized

    def _canonize_note(self, note: str) -> str | None:
        exact = super()._canonize_with_exact(note)
        if exact:
            return exact
        words = note.split()
        essential_words = [
            word
            for word in self._preprocess_note(words)
            if word not in self._meaningless_words
        ]
        without_adjectives, adjectives = self._divide_with_adjectives(essential_words)
        canonized = super().canonize(without_adjectives)
        if canonized:
            return canonized
        canonized = super().canonize(adjectives)
        if canonized:
            return canonized
        return super()._canonize_with_levenshtein(note)
