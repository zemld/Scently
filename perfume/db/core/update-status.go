package core

import "github.com/zemld/PerfumeRecommendationSystem/perfume/models"

type UpdateStatus struct {
	SuccessfulPerfumes []models.Perfume `json:"successful_perfumes"`
	FailedPerfumes     []models.Perfume `json:"failed_perfumes"`
	State              ProcessedState   `json:"state"`
}

func NewUpdateStatus(success bool) *UpdateStatus {
	status := UpdateStatus{
		SuccessfulPerfumes: []models.Perfume{},
		FailedPerfumes:     []models.Perfume{},
		State:              NewProcessedState(),
	}
	status.State.Success = success
	return &status
}
