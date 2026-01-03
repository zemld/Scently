package config

import "time"

const (
	HTTPClientTimeout          = 30 * time.Second
	HTTPClientMaxIdleConns     = 100
	HTTPClientMaxIdleConnsHost = 10
	HTTPClientIdleConnTimeout  = 90 * time.Second
)

const (
	AIFetcherTimeout         = 20 * time.Second
	PerfumeHubFetcherTimeout = 5 * time.Second
)

const (
	FamilyWeight          = 0.4
	NotesWeight           = 0.55
	TypeWeight            = 0.05
	UpperNotesWeight      = 0.2
	CoreNotesWeight       = 0.35
	BaseNotesWeight       = 0.45
	CharacteristicsWeight = 0.3
	TagsWeight            = 0.5
	OverlayWeight         = 0.2
)

const (
	ThreadsCount = 8
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
