package models

import "encoding/json"

type Perfume struct {
	Brand       string   `json:"brand"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Sex         string   `json:"sex"`
	Family      []string `json:"family"`
	UpperNotes  []string `json:"upper_notes"`
	MiddleNotes []string `json:"middle_notes"`
	BaseNotes   []string `json:"base_notes"`
	Link        string   `json:"link"`
	Volume      int      `json:"volume"`
	ImageUrl    string   `json:"image_url"`
}

type GluedPerfume struct {
	Brand      string            `json:"brand"`
	Name       string            `json:"name"`
	Properties PerfumeProperties `json:"properties"`
	Links      map[int]string    `json:"links"`
	ImageUrl   string            `json:"image_url"`
}

func NewGluedPerfume(p Perfume) GluedPerfume {
	return GluedPerfume{
		Brand:      p.Brand,
		Name:       p.Name,
		Properties: p.getProperties(),
		Links:      map[int]string{p.Volume: p.Link},
		ImageUrl:   p.ImageUrl,
	}
}

func (g GluedPerfume) Equal(other GluedPerfume) bool {
	return g.Brand == other.Brand && g.Name == other.Name
}

type PerfumeProperties struct {
	Type        string   `json:"type"`
	Sex         string   `json:"sex"`
	Family      []string `json:"family"`
	UpperNotes  []string `json:"upper_notes"`
	MiddleNotes []string `json:"middle_notes"`
	BaseNotes   []string `json:"base_notes"`
}

func (p Perfume) getProperties() PerfumeProperties {
	encodedPerfume, _ := json.Marshal(p)
	var props PerfumeProperties
	json.Unmarshal(encodedPerfume, &props)
	return props
}

type State struct {
	Success         bool `json:"success"`
	SuccessfulCount int  `json:"successful_count"`
	FailedCount     int  `json:"failed_count"`
}

type PerfumeResponse struct {
	Perfumes []Perfume `json:"perfumes"`
	State    State     `json:"state"`
}

type RankedPerfumeWithProps struct {
	Perfume GluedPerfume `json:"perfume"`
	Rank    int          `json:"rank"`
	Score   float64      `json:"similarity_score"`
}

type GluedPerfumeWithScore struct {
	GluedPerfume
	Score float64
}
