package models

import (
	"encoding/json"
)

type PerfumeProperties struct {
	Type        string   `json:"type"`
	Sex         string   `json:"sex"`
	Family      string   `json:"family"`
	UpperNotes  []string `json:"upper_notes"`
	MiddleNotes []string `json:"middle_notes"`
	BaseNotes   []string `json:"base_notes"`
}
type Perfume struct {
	Brand       string   `json:"brand"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Sex         string   `json:"sex"`
	Family      string   `json:"family"`
	UpperNotes  []string `json:"upper_notes"`
	MiddleNotes []string `json:"middle_notes"`
	BaseNotes   []string `json:"base_notes"`
	Link        string   `json:"link"`
	Volume      int      `json:"volume"`
}

type GluedPerfume struct {
	Brand      string            `json:"brand"`
	Name       string            `json:"name"`
	Properties PerfumeProperties `json:"properties"`
	Links      map[int]string    `json:"links"`
}

func NewGluedPerfume(p Perfume) GluedPerfume {
	return GluedPerfume{
		Brand:      p.Brand,
		Name:       p.Name,
		Properties: p.getProperties(),
		Links:      map[int]string{p.Volume: p.Link},
	}
}

func (p *Perfume) getProperties() PerfumeProperties {
	encodedPerfume, _ := json.Marshal(*p)
	var props PerfumeProperties
	json.Unmarshal(encodedPerfume, &props)
	return props
}
