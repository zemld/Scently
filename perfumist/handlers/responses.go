package handlers

import "github.com/zemld/PerfumeRecommendationSystem/perfumist/models"

type rankedPerfumeWithProps struct {
	Perfume models.GluedPerfume `json:"perfume"`
	Rank    int                 `json:"rank"`
}

type inputPerfume struct {
	Brand string `json:"brand"`
	Name  string `json:"name"`
	Ok    bool   `json:"ok"`
}

type SuggestResponse struct {
	Input     inputPerfume             `json:"input"`
	Suggested []rankedPerfumeWithProps `json:"suggested"`
	Success   bool                     `json:"success"`
}
