package core

type ProcessedState struct {
	Success         bool `json:"success"`
	SuccessfulCount int  `json:"successful_count"`
	FailedCount     int  `json:"failed_count"`
}

func NewProcessedState() ProcessedState {
	return ProcessedState{
		Success:         true,
		SuccessfulCount: 0,
		FailedCount:     0,
	}

}
