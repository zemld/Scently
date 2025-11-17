package config

import "time"

const (
	HTTPClientTimeout          = 30 * time.Second
	HTTPClientMaxIdleConns     = 100
	HTTPClientMaxIdleConnsHost = 10
	HTTPClientIdleConnTimeout  = 90 * time.Second
)

const (
	DBFetcherTimeout = 2 * time.Second
	AIFetcherTimeout = 20 * time.Second
)

const (
	FamilyWeight      = 0.4
	NotesWeight       = 0.55
	TypeWeight        = 0.05
	UpperNotesWeight  = 0.15
	MiddleNotesWeight = 0.45
	BaseNotesWeight   = 0.4
)

const (
	ThreadsCount = 5
	SuggestCount = 4
)

const (
	DefaultAISuggestURL   = "http://ai_advisor:8000/v1/advise"
	DefaultGetPerfumesURL = "http://perfume:8000/v1/perfumes/get"
)

const (
	PerfumeInternalTokenEnv = "PERFUME_INTERNAL_TOKEN"
	AISuggestURLEnv         = "AI_SUGGEST_URL"
	GetPerfumesURLEnv       = "GET_PERFUMES_URL"
)
