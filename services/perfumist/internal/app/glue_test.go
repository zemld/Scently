package app

import (
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models"
)

func TestGlue_MergesLinksByBrandName(t *testing.T) {
	t.Parallel()

	perfumes := []models.Perfume{
		{Brand: "A", Name: "X", Link: "l1", Volume: 50},
		{Brand: "A", Name: "X", Link: "l2", Volume: 100},
		{Brand: "B", Name: "Y", Link: "l3", Volume: 30},
	}

	glued := Glue(perfumes)

	if len(glued) != 2 {
		t.Fatalf("expected 2 glued perfumes, got %d", len(glued))
	}

	var ax *models.GluedPerfume
	var by *models.GluedPerfume
	for i := range glued {
		g := glued[i]
		if g.Brand == "A" && g.Name == "X" {
			ax = &g
		}
		if g.Brand == "B" && g.Name == "Y" {
			by = &g
		}
	}

	if ax == nil || by == nil {
		t.Fatalf("expected both A+X and B+Y in glued result")
	}

	if len(ax.Links) != 2 || ax.Links[50] != "l1" || ax.Links[100] != "l2" {
		t.Fatalf("unexpected links for A+X: %+v", ax.Links)
	}

	if len(by.Links) != 1 || by.Links[30] != "l3" {
		t.Fatalf("unexpected links for B+Y: %+v", by.Links)
	}
}

func TestGetKey(t *testing.T) {
	t.Parallel()

	p := models.Perfume{Brand: "A", Name: "X"}
	if got := getKey(p); got != "AX" {
		t.Fatalf("getKey expected AX, got %q", got)
	}
}

func TestFetchGluedPerfumesFromMap_Size(t *testing.T) {
	t.Parallel()

	m := map[string]models.GluedPerfume{
		"AX": {Brand: "A", Name: "X", Links: map[int]string{50: "l1"}},
		"BY": {Brand: "B", Name: "Y", Links: map[int]string{30: "l2"}},
	}
	res := fetchGluedPerfumesFromMap(m)
	if len(res) != 2 {
		t.Fatalf("expected 2, got %d", len(res))
	}
}
