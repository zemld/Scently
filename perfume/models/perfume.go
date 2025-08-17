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
	Volumes     []int    `json:"volumes"`
	Links       []string `json:"links"`
}

func (p *Perfume) Unpack() []any {
	fields := make([]any, 0)
	v := reflect.ValueOf(*p)
	for i := 0; i < v.NumField(); i++ {
		fields = append(fields, v.Field(i).Interface())
	}
	return fields
}
