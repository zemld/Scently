def _initialize_distance_matrix(a: str, b: str) -> list[list[int]]:
    matrix = [[0] * (len(b) + 1) for _ in range(len(a) + 1)]
    for i in range(len(a) + 1):
        matrix[i][0] = i
    for j in range(len(b) + 1):
        matrix[0][j] = j
    return matrix


def get_levenshtein_distance(a: str, b: str) -> int:
    if a == b:
        return 0
    if not a:
        return len(b)
    if not b:
        return len(a)

    dist_matrix = _initialize_distance_matrix(a, b)
    for i in range(1, len(a) + 1):
        for j in range(1, len(b) + 1):
            cost = 0 if a[i - 1] == b[j - 1] else 1
            dist_matrix[i][j] = min(
                dist_matrix[i - 1][j] + 1,
                dist_matrix[i][j - 1] + 1,
                dist_matrix[i - 1][j - 1] + cost,
            )

    return dist_matrix[len(a)][len(b)]
