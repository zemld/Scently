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
	Volume      int      `json:"volume"`
	Link        string   `json:"link"`
}

type GluedPerfume struct {
	Perfume
	Volumes []int    `json:"volumes"`
	Links   []string `json:"links"`
}
