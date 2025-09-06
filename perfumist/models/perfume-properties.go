package models

import "encoding/json"

type PerfumeProperties struct {
	Type        string   `json:"type"`
	Sex         string   `json:"sex"`
	Family      []string `json:"family"`
	UpperNotes  []string `json:"upper_notes"`
	MiddleNotes []string `json:"middle_notes"`
	BaseNotes   []string `json:"base_notes"`
}

func (p *Perfume) getProperties() PerfumeProperties {
	encodedPerfume, _ := json.Marshal(*p)
	var props PerfumeProperties
	json.Unmarshal(encodedPerfume, &props)
	return props
}
