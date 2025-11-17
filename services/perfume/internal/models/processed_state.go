package models

type ProcessedState struct {
	SuccessfulCount int   `json:"successful_count"`
	FailedCount     int   `json:"failed_count"`
	Error           error `json:"-"`
}

func NewProcessedState() ProcessedState {
	return ProcessedState{
		SuccessfulCount: 0,
		FailedCount:     0,
		Error:           nil,
	}
}
