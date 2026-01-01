package config

import "time"

const (
	HTTPClientTimeout          = 30 * time.Second
	HTTPClientMaxIdleConns     = 100
	HTTPClientMaxIdleConnsHost = 10
	HTTPClientIdleConnTimeout  = 90 * time.Second
)

const (
	DBFetcherTimeout         = 2 * time.Second
	AIFetcherTimeout         = 20 * time.Second
	PerfumeHubFetcherTimeout = 1 * time.Second
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
	DefaultGetPerfumesURL = "http://perfume:8000/v1/perfumes/get"
)

const (
	PerfumeInternalTokenEnv = "PERFUME_INTERNAL_TOKEN"
	GetPerfumesURLEnv       = "GET_PERFUMES_URL"
	PerfumeHubURLEnv        = "PERFUME_HUB_URL"
)
