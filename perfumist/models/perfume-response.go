package models

type State struct {
	Success         bool `json:"success"`
	SuccessfulCount int  `json:"successful_count"`
	FailedCount     int  `json:"failed_count"`
}

type PerfumeResponse struct {
	Perfumes []Perfume `json:"perfumes"`
	State    State     `json:"state"`
}
