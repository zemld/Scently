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
			constants.Select + " WHERE sexes.sex = 'unisex'", nil},
		{"brand only", NewSelectParameters().WithBrand("Chanel"),
			constants.Select + " WHERE perfumes.brand = $1 AND sexes.sex = 'unisex'", []any{"Chanel"}},
		{"name only", NewSelectParameters().WithName("No.5"),
			constants.Select + " WHERE perfumes.name = $1 AND sexes.sex = 'unisex'", []any{"No.5"}},
		{"sex unisex", NewSelectParameters().WithSex("unisex"),
			constants.Select + " WHERE sexes.sex = 'unisex'", nil},
		{"sex female", NewSelectParameters().WithSex("female"),
			constants.Select + " WHERE (sexes.sex = $1 OR sexes.sex = 'unisex')", []any{"female"}},
		{"sex male", NewSelectParameters().WithSex("male"),
			constants.Select + " WHERE (sexes.sex = $1 OR sexes.sex = 'unisex')", []any{"male"}},
		{"brand and name", NewSelectParameters().WithBrand("Dior").WithName("Sauvage"),
			constants.Select + " WHERE perfumes.brand = $1 AND perfumes.name = $2 AND sexes.sex = 'unisex'", []any{"Dior", "Sauvage"}},
		{"brand and sex", NewSelectParameters().WithBrand("Chanel").WithSex("female"),
			constants.Select + " WHERE perfumes.brand = $1 AND (sexes.sex = $2 OR sexes.sex = 'unisex')", []any{"Chanel", "female"}},
		{"name and sex unisex", NewSelectParameters().WithName("No.5").WithSex("unisex"),
			constants.Select + " WHERE perfumes.name = $1 AND sexes.sex = 'unisex'", []any{"No.5"}},
		{"name and sex female", NewSelectParameters().WithName("No.5").WithSex("female"),
			constants.Select + " WHERE perfumes.name = $1 AND (sexes.sex = $2 OR sexes.sex = 'unisex')", []any{"No.5", "female"}},
		{"brand, name and sex", NewSelectParameters().WithBrand("Dior").WithName("Sauvage").WithSex("male"),
			constants.Select + " WHERE perfumes.brand = $1 AND perfumes.name = $2 AND (sexes.sex = $3 OR sexes.sex = 'unisex')", []any{"Dior", "Sauvage", "male"}},
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
