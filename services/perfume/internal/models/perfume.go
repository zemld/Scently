package models

type Perfume struct {
	Brand      string            `json:"brand"`
	Name       string            `json:"name"`
	Sex        string            `json:"sex"`
	Properties PerfumeProperties `json:"properties"`
	Shops      []ShopInfo        `json:"shops"`
}

type PerfumeProperties struct {
	Type       string   `json:"type"`
	Family     []string `json:"family"`
	UpperNotes []string `json:"upper_notes"`
	CoreNotes  []string `json:"core_notes"`
	BaseNotes  []string `json:"base_notes"`
}

type ShopInfo struct {
	ShopName string           `json:"shop_name"`
	Domain   string           `json:"domain"`
	ImageUrl string           `json:"image_url"`
	Variants []PerfumeVariant `json:"variants"`
}

type PerfumeVariant struct {
	Volume int    `json:"volume"`
	Link   string `json:"link"`
	Price  int    `json:"price"`
}
