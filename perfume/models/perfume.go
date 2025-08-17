package models

import "reflect"

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

func (p *GluedPerfume) Unpack() []any {
	fields := make([]any, 0)
	v := reflect.ValueOf(*p)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Type() == reflect.TypeOf(p.Perfume) {
			fields = p.Perfume.unpack(fields)
		} else {
			fields = append(fields, v.Field(i).Interface())
		}
	}
	return fields
}

func (p *Perfume) unpack(fields []any) []any {
	v := reflect.ValueOf(*p)
	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Name != "Volume" && v.Type().Field(i).Name != "Link" {
			fields = append(fields, v.Field(i).Interface())
		}
	}
	return fields
}
