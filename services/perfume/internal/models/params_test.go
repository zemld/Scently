package models

import (
	"reflect"
	"strings"
	"testing"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/constants"
)

func TestSelectParameters_getQuery(t *testing.T) {
	baseQuery := strings.TrimSpace(constants.Select)
	groupByIndex := strings.LastIndex(baseQuery, "GROUP BY")

	tests := []struct {
		name     string
		p        *SelectParameters
		want     string
		wantArgs []any
	}{
		{"no filters", NewSelectParameters(),
			baseQuery[:groupByIndex] + " WHERE s.sex = 'unisex' " + baseQuery[groupByIndex:], nil},
		{"brand only", NewSelectParameters().WithBrand("Chanel"),
			baseQuery[:groupByIndex] + " WHERE pb.brand = $1 AND s.sex = 'unisex' " + baseQuery[groupByIndex:], []any{"Chanel"}},
		{"name only", NewSelectParameters().WithName("No.5"),
			baseQuery[:groupByIndex] + " WHERE pb.name = $1 AND s.sex = 'unisex' " + baseQuery[groupByIndex:], []any{"No.5"}},
		{"sex unisex", NewSelectParameters().WithSex("unisex"),
			baseQuery[:groupByIndex] + " WHERE s.sex = 'unisex' " + baseQuery[groupByIndex:], nil},
		{"sex female", NewSelectParameters().WithSex("female"),
			baseQuery[:groupByIndex] + " WHERE (s.sex = $1 OR s.sex = 'unisex') " + baseQuery[groupByIndex:], []any{"female"}},
		{"sex male", NewSelectParameters().WithSex("male"),
			baseQuery[:groupByIndex] + " WHERE (s.sex = $1 OR s.sex = 'unisex') " + baseQuery[groupByIndex:], []any{"male"}},
		{"brand and name", NewSelectParameters().WithBrand("Dior").WithName("Sauvage"),
			baseQuery[:groupByIndex] + " WHERE pb.brand = $1 AND pb.name = $2 AND s.sex = 'unisex' " + baseQuery[groupByIndex:], []any{"Dior", "Sauvage"}},
		{"brand and sex", NewSelectParameters().WithBrand("Chanel").WithSex("female"),
			baseQuery[:groupByIndex] + " WHERE pb.brand = $1 AND (s.sex = $2 OR s.sex = 'unisex') " + baseQuery[groupByIndex:], []any{"Chanel", "female"}},
		{"name and sex unisex", NewSelectParameters().WithName("No.5").WithSex("unisex"),
			baseQuery[:groupByIndex] + " WHERE pb.name = $1 AND s.sex = 'unisex' " + baseQuery[groupByIndex:], []any{"No.5"}},
		{"name and sex female", NewSelectParameters().WithName("No.5").WithSex("female"),
			baseQuery[:groupByIndex] + " WHERE pb.name = $1 AND (s.sex = $2 OR s.sex = 'unisex') " + baseQuery[groupByIndex:], []any{"No.5", "female"}},
		{"brand, name and sex", NewSelectParameters().WithBrand("Dior").WithName("Sauvage").WithSex("male"),
			baseQuery[:groupByIndex] + " WHERE pb.brand = $1 AND pb.name = $2 AND (s.sex = $3 OR s.sex = 'unisex') " + baseQuery[groupByIndex:], []any{"Dior", "Sauvage", "male"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetQuery(); got != tt.want {
				t.Fatalf("getQuery() = %q, want %q", got, tt.want)
			}
			if gotArgs := tt.p.Unpack(); !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Fatalf("unpack() = %#v, want %#v", gotArgs, tt.wantArgs)
			}
		})
	}
}
