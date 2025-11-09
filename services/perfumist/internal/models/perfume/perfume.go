package perfume

type Perfume struct {
	Brand      string     `json:"brand"`
	Name       string     `json:"name"`
	Sex        string     `json:"sex"`
	ImageUrl   string     `json:"image_url"`
	Properties Properties `json:"properties"`
	Shops      []ShopInfo `json:"shops"`
}

type Properties struct {
	Type       string   `json:"perfume_type"`
	Family     []string `json:"family"`
	UpperNotes []string `json:"upper_notes"`
	CoreNotes  []string `json:"core_notes"`
	BaseNotes  []string `json:"base_notes"`
}

type ShopInfo struct {
	ShopName string    `json:"shop_name"`
	Domain   string    `json:"domain"`
	Variants []Variant `json:"variants"`
}

type Variant struct {
	Volume int    `json:"volume"`
	Link   string `json:"link"`
	Price  int    `json:"price"`
}

func (p Perfume) Equal(other Perfume) bool {
	return p.Brand == other.Brand && p.Name == other.Name && p.Sex == other.Sex
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

type Ranked struct {
	Perfume Perfume `json:"perfume"`
	Rank    int     `json:"rank,omitempty"`
	Score   float64 `json:"similarity_score,omitempty"`
}

type WithScore struct {
	Perfume
	Score float64
}
