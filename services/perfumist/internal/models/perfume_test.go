package models

import "testing"

func TestPerfume_Equal(t *testing.T) {
	t.Parallel()

	p := Perfume{
		Brand:    "A",
		Name:     "X",
		Sex:      "male",
		ImageUrl: "url",
		Properties: PerfumeProperties{
			Type:       "edt",
			Family:     []string{"woody"},
			UpperNotes: []string{"bergamot"},
			CoreNotes:  []string{"lavender"},
			BaseNotes:  []string{"cedar"},
		},
	}

	if !p.Equal(p) {
		t.Fatalf("Equal should be true for same perfume")
	}

	// Test different name
	other := Perfume{Brand: "A", Name: "Y", Sex: "male"}
	if p.Equal(other) {
		t.Fatalf("Equal should be false for different name")
	}

	// Test different sex
	otherSex := Perfume{Brand: "A", Name: "X", Sex: "female"}
	if p.Equal(otherSex) {
		t.Fatalf("Equal should be false for different sex")
	}

	// Test different brand
	otherBrand := Perfume{Brand: "B", Name: "X", Sex: "male"}
	if p.Equal(otherBrand) {
		t.Fatalf("Equal should be false for different brand")
	}
}

func TestPerfumeProperties(t *testing.T) {
	t.Parallel()

	props := PerfumeProperties{
		Type:       "edt",
		Family:     []string{"woody"},
		UpperNotes: []string{"bergamot"},
		CoreNotes:  []string{"lavender"},
		BaseNotes:  []string{"cedar"},
	}

	if props.Type != "edt" {
		t.Fatalf("type should be edt, got %s", props.Type)
	}
	if len(props.Family) != 1 || props.Family[0] != "woody" {
		t.Fatalf("family should contain woody, got %+v", props.Family)
	}
	if len(props.UpperNotes) != 1 || props.UpperNotes[0] != "bergamot" {
		t.Fatalf("upper notes should contain bergamot, got %+v", props.UpperNotes)
	}
	if len(props.CoreNotes) != 1 || props.CoreNotes[0] != "lavender" {
		t.Fatalf("core notes should contain lavender, got %+v", props.CoreNotes)
	}
	if len(props.BaseNotes) != 1 || props.BaseNotes[0] != "cedar" {
		t.Fatalf("base notes should contain cedar, got %+v", props.BaseNotes)
	}
}
