package models

type Perfume struct {
	Brand      string            `json:"brand"`
	Name       string            `json:"name"`
	Sex        string            `json:"sex"`
	ImageUrl   string            `json:"image_url"`
	Properties PerfumeProperties `json:"properties"`
	Shops      []ShopInfo        `json:"shops"`
}

type PerfumeProperties struct {
	Type       string   `json:"perfume_type"`
	Family     []string `json:"family"`
	UpperNotes []string `json:"upper_notes"`
	CoreNotes  []string `json:"core_notes"`
	BaseNotes  []string `json:"base_notes"`
}

type ShopInfo struct {
	ShopName string           `json:"shop_name"`
	Domain   string           `json:"domain"`
	ImageUrl string           `json:"image_url,omitempty"`
	Variants []PerfumeVariant `json:"variants"`
}

type PerfumeVariant struct {
	Volume int    `json:"volume"`
	Link   string `json:"link"`
	Price  int    `json:"price"`
}

type CanonizedPerfume struct {
	Brand string
	Name  string
}

func (p Perfume) Canonize() CanonizedPerfume {
	return CanonizedPerfume{
		Brand: canonize(p.Brand),
		Name:  canonize(p.Name),
	}
}
