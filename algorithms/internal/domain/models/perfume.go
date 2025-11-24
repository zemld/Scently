package models

type Perfume struct {
	Brand      string     `json:"brand"`
	Name       string     `json:"name"`
	Sex        string     `json:"sex"`
	ImageUrl   string     `json:"image_url"`
	Properties Properties `json:"properties"`
	Shops      []ShopInfo `json:"shops"`
}

type Properties struct {
	Type   string   `json:"perfume_type"`
	Family []string `json:"family"`

	UpperNotes []string `json:"upper_notes"`
	CoreNotes  []string `json:"core_notes"`
	BaseNotes  []string `json:"base_notes"`

	EnrichedUpperNotes []EnrichedNote `json:"-"`
	EnrichedCoreNotes  []EnrichedNote `json:"-"`
	EnrichedBaseNotes  []EnrichedNote `json:"-"`

	UpperTags map[string]int `json:"upper_tags,omitempty"`
	CoreTags  map[string]int `json:"core_tags,omitempty"`
	BaseTags  map[string]int `json:"base_tags,omitempty"`

	Tags map[string]int `json:"tags,omitempty"`

	UpperCharacteristics map[string]float64 `json:"upper_characteristics,omitempty"`
	CoreCharacteristics  map[string]float64 `json:"core_characteristics,omitempty"`
	BaseCharacteristics  map[string]float64 `json:"base_characteristics,omitempty"`

	Characteristics map[string]float64 `json:"characteristics,omitempty"`
}

type EnrichedNote struct {
	Name            string
	Tags            map[string]int
	Characteristics map[string]float64
}

type ShopInfo struct {
	ShopName string    `json:"shop_name"`
	Domain   string    `json:"domain"`
	ImageUrl string    `json:"image_url,omitempty"`
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

func UniteTags(notes []EnrichedNote) map[string]int {
	united := make(map[string]int)

	for _, note := range notes {
		for tag, count := range note.Tags {
			united[tag] += count
		}
	}

	return united
}

func UniteCharacteristics(notes []EnrichedNote) map[string]float64 {
	united := make(map[string]float64)

	for _, note := range notes {
		for characteristic, value := range note.Characteristics {
			united[characteristic] += value
		}
	}

	for characteristic, value := range united {
		united[characteristic] = value / float64(len(notes))
	}

	return united
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
