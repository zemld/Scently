package models

type RunResults struct {
	ID   string                `json:"id"`
	Runs []SingleSuggestResult `json:"runs"`
}

type SingleSuggestResult struct {
	Input       Perfume  `json:"input"`
	Suggestions []Ranked `json:"suggestions"`
}

type ShortenRunResults struct {
	ID   string                       `json:"id"`
	Runs []ShortenSingleSuggestResult `json:"runs"`
}

type ShortenSingleSuggestResult struct {
	Input       ShortenPerfume   `json:"input"`
	Suggestions []ShortenPerfume `json:"suggestions"`
}

type ShortenPerfume struct {
	Brand string `json:"brand"`
	Name  string `json:"name"`
	Sex   string `json:"sex"`
}

func NewShortenRunResults(runResults RunResults) *ShortenRunResults {
	result := &ShortenRunResults{
		ID:   runResults.ID,
		Runs: make([]ShortenSingleSuggestResult, 0, len(runResults.Runs)),
	}
	for _, runResult := range runResults.Runs {
		result.Runs = append(result.Runs, *NewShortenSingleSuggestResult(runResult))
	}
	return result
}

func NewShortenSingleSuggestResult(singleSuggestResult SingleSuggestResult) *ShortenSingleSuggestResult {
	result := &ShortenSingleSuggestResult{
		Input:       *NewShortenPerfume(singleSuggestResult.Input),
		Suggestions: make([]ShortenPerfume, 0, len(singleSuggestResult.Suggestions)),
	}
	for _, suggestion := range singleSuggestResult.Suggestions {
		result.Suggestions = append(result.Suggestions, *NewShortenPerfume(suggestion.Perfume))
	}
	return result
}

func NewShortenPerfume(perfume Perfume) *ShortenPerfume {
	return &ShortenPerfume{
		Brand: perfume.Brand,
		Name:  perfume.Name,
		Sex:   perfume.Sex,
	}
}
