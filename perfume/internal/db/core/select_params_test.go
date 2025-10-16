package core

import (
	"reflect"
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/constants"
)

func TestSelectParameters_getQuery(t *testing.T) {
	tests := []struct {
		name     string
		p        *SelectParameters
		want     string
		wantArgs []any
	}{
		{"no filters", NewSelectParameters(),
			constants.Select, nil},
		{"brand only", NewSelectParameters().WithBrand("Chanel"),
			constants.Select + " WHERE perfumes.brand = $1", []any{"Chanel"}},
		{"name only", NewSelectParameters().WithName("No.5"),
			constants.Select + " WHERE perfumes.name = $1", []any{"No.5"}},
		{"brand and name", NewSelectParameters().WithBrand("Dior").WithName("Sauvage"),
			constants.Select + " WHERE perfumes.brand = $1 AND perfumes.name = $2", []any{"Dior", "Sauvage"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.getQuery(); got != tt.want {
				t.Fatalf("getQuery() = %q, want %q", got, tt.want)
			}
			if gotArgs := tt.p.unpack(); !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Fatalf("unpack() = %#v, want %#v", gotArgs, tt.wantArgs)
			}
		})
	}
}
