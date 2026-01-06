package models

import (
	"strings"
	"unicode"
)

const (
	Male   Sex = "male"
	Female Sex = "female"
	Unisex Sex = "unisex"
)

type Perfume struct {
	Brand      string     `json:"brand"`
	Name       string     `json:"name"`
	Sex        Sex        `json:"sex"`
	ImageUrl   string     `json:"image_url"`
	Properties Properties `json:"properties"`
	Shops      []ShopInfo `json:"shops"`
}

type Sex string

type Properties struct {
	Type       string   `json:"perfume_type"`
	Family     []string `json:"family"`
	UpperNotes []string `json:"upper_notes"`
	CoreNotes  []string `json:"core_notes"`
	BaseNotes  []string `json:"base_notes"`

	EnrichedUpperNotes []EnrichedNote `json:"enriched_upper_notes,omitempty"`
	EnrichedCoreNotes  []EnrichedNote `json:"enriched_core_notes,omitempty"`
	EnrichedBaseNotes  []EnrichedNote `json:"enriched_base_notes,omitempty"`

	Tags map[string]int `json:"tags,omitempty"`

	UpperCharacteristics map[string]float64 `json:"upper_characteristics,omitempty"`
	CoreCharacteristics  map[string]float64 `json:"core_characteristics,omitempty"`
	BaseCharacteristics  map[string]float64 `json:"base_characteristics,omitempty"`
}

type EnrichedNote struct {
	Name            string               `json:"name"`
	Tags            []string             `json:"tags"`
	Characteristics []NoteCharacteristic `json:"characteristics"`
}

type NoteCharacteristic struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
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

type CanonizedPerfume struct {
	Brand string
	Name  string
	Sex   Sex
}

func (p Perfume) Canonize() CanonizedPerfume {
	return CanonizedPerfume{
		Brand: canonize(p.Brand),
		Name:  canonize(p.Name),
		Sex:   p.Sex,
	}
}

func canonize(s string) string {
	canonized := strings.Builder{}

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			continue
		}
		canonized.WriteRune(unicode.ToLower(r))
	}
	return canonized.String()
}

type Ranked struct {
	Perfume Perfume `json:"perfume"`
	Rank    int     `json:"rank,omitzero"`
	Score   float64 `json:"similarity_score"`
}
