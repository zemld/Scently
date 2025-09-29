package similarity

func getTypeSimilarityScore(first string, second string) float64 {
	if first == second {
		return 1
	}
	return 0
}
