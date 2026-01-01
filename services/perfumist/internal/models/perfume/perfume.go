package perfume

import "github.com/zemld/Scently/models"

type State struct {
	SuccessfulCount int `json:"successful_count"`
	FailedCount     int `json:"failed_count"`
}

type PerfumeResponse struct {
	Perfumes []models.Perfume `json:"perfumes"`
	State    State            `json:"state"`
}
