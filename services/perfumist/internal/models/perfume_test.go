package models

import "testing"

func TestNewGluedPerfume_AndEqual(t *testing.T) {
	t.Parallel()

	p := Perfume{
		Brand:       "A",
		Name:        "X",
		Type:        "edt",
		Sex:         "male",
		Family:      []string{"woody"},
		UpperNotes:  []string{"bergamot"},
		MiddleNotes: []string{"lavender"},
		BaseNotes:   []string{"cedar"},
		Link:        "link",
		Volume:      50,
		ImageUrl:    "url",
	}

	g := NewGluedPerfume(p)
	if g.Brand != p.Brand || g.Name != p.Name {
		t.Fatalf("brand/name not copied")
	}
	if g.Properties.Type != p.Type {
		t.Fatalf("properties not copied correctly: %+v", g.Properties)
	}
	if len(g.Links) != 1 || g.Links[50] != "link" {
		t.Fatalf("links not initialized correctly: %+v", g.Links)
	}
	if g.ImageUrl != p.ImageUrl {
		t.Fatalf("image url not copied correctly: %+v", g.ImageUrl)
	}

	if !g.Equal(NewGluedPerfume(p)) {
		t.Fatalf("Equal should be true for same brand+name+sex")
	}

	// Test different name
	other := GluedPerfume{Brand: "A", Name: "Y", Sex: "male"}
	if g.Equal(other) {
		t.Fatalf("Equal should be false for different name")
	}

	// Test different sex
	otherSex := GluedPerfume{Brand: "A", Name: "X", Sex: "female"}
	if g.Equal(otherSex) {
		t.Fatalf("Equal should be false for different sex")
	}

	// Test different brand
	otherBrand := GluedPerfume{Brand: "B", Name: "X", Sex: "male"}
	if g.Equal(otherBrand) {
		t.Fatalf("Equal should be false for different brand")
	}
}

func TestGetProperties(t *testing.T) {
	t.Parallel()

	p := Perfume{
		Brand:       "A",
		Name:        "X",
		Type:        "edt",
		Sex:         "male",
		Family:      []string{"woody"},
		UpperNotes:  []string{"bergamot"},
		MiddleNotes: []string{"lavender"},
		BaseNotes:   []string{"cedar"},
		Link:        "link",
		Volume:      50,
	}
	props := p.getProperties()
	if props.Type != p.Type || len(props.Family) != 1 {
		t.Fatalf("properties mapping incorrect: %+v", props)
	}
}
