package models

import (
	"reflect"
	"testing"
)

func TestPerfumeUnpackPropertiesAndLinkedFields(t *testing.T) {
	p := Perfume{
		Brand:       "BrandA",
		Name:        "NameX",
		Type:        "EDP",
		Sex:         "unisex",
		Family:      []string{"citrus", "woody"},
		UpperNotes:  []string{"bergamot"},
		MiddleNotes: []string{"rose"},
		BaseNotes:   []string{"musk"},
		Link:        "http://example.com",
		Volume:      100,
	}

	gotProps := p.UnpackProperties()
	gotLinked := p.UnpackLinkedFields()

	wantProps := []any{
		p.Brand,
		p.Name,
		p.Type,
		p.Sex,
		p.Family,
		p.UpperNotes,
		p.MiddleNotes,
		p.BaseNotes,
	}
	if !reflect.DeepEqual(gotProps, wantProps) {
		t.Fatalf("UnpackProperties() mismatch:\n got: %#v\nwant: %#v", gotProps, wantProps)
	}

	wantLinked := []any{
		p.Brand,
		p.Name,
		p.Link,
		p.Volume,
	}
	if !reflect.DeepEqual(gotLinked, wantLinked) {
		t.Fatalf("UnpackLinkedFields() mismatch:\n got: %#v\nwant: %#v", gotLinked, wantLinked)
	}
}
