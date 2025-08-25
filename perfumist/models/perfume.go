package models

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
	Perfume
	Links map[int]string `json:"links"`
}

func NewGluedPerfume(p Perfume) GluedPerfume {
	return GluedPerfume{Perfume: p, Links: map[int]string{p.Volume: p.Link}}
}
