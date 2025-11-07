package core

import "github.com/zemld/PerfumeRecommendationSystem/perfume/internal/models"

const UpdateParametersContextKey contextKey = "update_parameters"

type UpdateParameters struct {
	// TODO: Remove IsTruncate
	IsTruncate       bool
	Perfumes         []models.Perfume         `json:"perfumes"`
	UpgradedPerfumes []models.UpgradedPerfume `json:"upgraded_perfumes"`
}

func NewUpdateParameters() *UpdateParameters {
	return &UpdateParameters{
		IsTruncate: false,
	}
}

func (p *UpdateParameters) WithTruncate() *UpdateParameters {
	p.IsTruncate = true
	return p
}

func (p *UpdateParameters) WithPerfumes(perfumes []models.Perfume) *UpdateParameters {
	p.Perfumes = perfumes
	return p
}
