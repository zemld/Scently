package models

import "reflect"

type Perfume struct {
	Brand       string   `json:"brand"`
	Name        string   `json:"name"`
	Sex         string   `json:"sex"`
	Type        string   `json:"type"`
	Family      []string `json:"family"`
	UpperNotes  []string `json:"upper_notes"`
	MiddleNotes []string `json:"middle_notes"`
	BaseNotes   []string `json:"base_notes"`
	Link        string   `json:"link"`
	Volume      int      `json:"volume"`
	ImageUrl    string   `json:"image_url"`
}

func (p *Perfume) UnpackProperties() []any {
	fields := make([]any, 0)
	v := reflect.ValueOf(*p)
	for i := 0; i < v.NumField(); i++ {
		if isPropertyOrKey(v.Type().Field(i)) {
			fields = append(fields, v.Field(i).Interface())
		}
	}
	return fields
}

func (p *Perfume) UnpackLinkedFields() []any {
	fields := make([]any, 0)
	v := reflect.ValueOf(*p)
	for i := 0; i < v.NumField(); i++ {
		if isNotProperty(v.Type().Field(i)) {
			fields = append(fields, v.Field(i).Interface())
		}
	}
	return fields
}

func isPropertyOrKey(field reflect.StructField) bool {
	return isProperty(field) || field.Name == "Sex" || field.Name == "Brand" || field.Name == "Name"
}

func isProperty(field reflect.StructField) bool {
	return !isNotProperty(field)
}

func isNotProperty(field reflect.StructField) bool {
	return field.Name == "Link" || field.Name == "Volume" || field.Name == "Sex" || field.Name == "Brand" || field.Name == "Name"
}

type UpgradedPerfume struct {
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
