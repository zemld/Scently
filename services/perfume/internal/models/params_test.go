package models

import (
	"reflect"
	"strings"
	"testing"

	queries "github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/query"
)

func TestSelectParameters_GetChoosingPerfumesQuery(t *testing.T) {
	baseQuery := strings.TrimSpace(queries.SelectPerfumesBaseInfo)

	tests := []struct {
		name string
		p    *SelectParameters
		want string
	}{
		{
			"no filters - defaults to unisex",
			NewSelectParameters(),
			baseQuery + " WHERE s.sex = 'unisex'",
		},
		{
			"brand only",
			NewSelectParameters().WithBrand("Chanel"),
			baseQuery + " WHERE pb.brand = $1 AND s.sex = 'unisex'",
		},
		{
			"name only",
			NewSelectParameters().WithName("No.5"),
			baseQuery + " WHERE pb.name = $1 AND s.sex = 'unisex'",
		},
		{
			"sex unisex",
			NewSelectParameters().WithSex("unisex"),
			baseQuery + " WHERE s.sex = 'unisex'",
		},
		{
			"sex female",
			NewSelectParameters().WithSex("female"),
			baseQuery + " WHERE (s.sex = $1 OR s.sex = 'unisex')",
		},
		{
			"sex male",
			NewSelectParameters().WithSex("male"),
			baseQuery + " WHERE (s.sex = $1 OR s.sex = 'unisex')",
		},
		{
			"brand and name",
			NewSelectParameters().WithBrand("Dior").WithName("Sauvage"),
			baseQuery + " WHERE pb.brand = $1 AND pb.name = $2 AND s.sex = 'unisex'",
		},
		{
			"brand and sex female",
			NewSelectParameters().WithBrand("Chanel").WithSex("female"),
			baseQuery + " WHERE pb.brand = $1 AND (s.sex = $2 OR s.sex = 'unisex')",
		},
		{
			"brand and sex unisex",
			NewSelectParameters().WithBrand("Chanel").WithSex("unisex"),
			baseQuery + " WHERE pb.brand = $1 AND s.sex = 'unisex'",
		},
		{
			"name and sex unisex",
			NewSelectParameters().WithName("No.5").WithSex("unisex"),
			baseQuery + " WHERE pb.name = $1 AND s.sex = 'unisex'",
		},
		{
			"name and sex female",
			NewSelectParameters().WithName("No.5").WithSex("female"),
			baseQuery + " WHERE pb.name = $1 AND (s.sex = $2 OR s.sex = 'unisex')",
		},
		{
			"brand, name and sex male",
			NewSelectParameters().WithBrand("Dior").WithName("Sauvage").WithSex("male"),
			baseQuery + " WHERE pb.brand = $1 AND pb.name = $2 AND (s.sex = $3 OR s.sex = 'unisex')",
		},
		{
			"brand, name and sex unisex",
			NewSelectParameters().WithBrand("Dior").WithName("Sauvage").WithSex("unisex"),
			baseQuery + " WHERE pb.brand = $1 AND pb.name = $2 AND s.sex = 'unisex'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.GetChoosingPerfumesQuery()
			if got != tt.want {
				t.Errorf("GetChoosingPerfumesQuery() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSelectParameters_Unpack(t *testing.T) {
	tests := []struct {
		name     string
		p        *SelectParameters
		wantArgs []any
	}{
		{
			"no filters",
			NewSelectParameters(),
			nil,
		},
		{
			"brand only",
			NewSelectParameters().WithBrand("Chanel"),
			[]any{"Chanel"},
		},
		{
			"name only",
			NewSelectParameters().WithName("No.5"),
			[]any{"No.5"},
		},
		{
			"sex unisex - no args",
			NewSelectParameters().WithSex("unisex"),
			nil,
		},
		{
			"sex female",
			NewSelectParameters().WithSex("female"),
			[]any{"female"},
		},
		{
			"sex male",
			NewSelectParameters().WithSex("male"),
			[]any{"male"},
		},
		{
			"brand and name",
			NewSelectParameters().WithBrand("Dior").WithName("Sauvage"),
			[]any{"Dior", "Sauvage"},
		},
		{
			"brand and sex female",
			NewSelectParameters().WithBrand("Chanel").WithSex("female"),
			[]any{"Chanel", "female"},
		},
		{
			"brand and sex unisex - no sex arg",
			NewSelectParameters().WithBrand("Chanel").WithSex("unisex"),
			[]any{"Chanel"},
		},
		{
			"name and sex female",
			NewSelectParameters().WithName("No.5").WithSex("female"),
			[]any{"No.5", "female"},
		},
		{
			"brand, name and sex male",
			NewSelectParameters().WithBrand("Dior").WithName("Sauvage").WithSex("male"),
			[]any{"Dior", "Sauvage", "male"},
		},
		{
			"brand, name and sex unisex - no sex arg",
			NewSelectParameters().WithBrand("Dior").WithName("Sauvage").WithSex("unisex"),
			[]any{"Dior", "Sauvage"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArgs := tt.p.Unpack()
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("Unpack() = %#v, want %#v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestSelectParameters_GetQuery(t *testing.T) {
	tests := []struct {
		name string
		p    *SelectParameters
		want string
	}{
		{
			"no filters - defaults to unisex",
			NewSelectParameters(),
			func() string {
				choosingQuery := strings.TrimSpace(queries.SelectPerfumesBaseInfo) + " WHERE s.sex = 'unisex'"
				withClause := strings.Replace(queries.WithSelect, "%s", choosingQuery, 1)
				return withClause + queries.EnrichSelectedPerfumes
			}(),
		},
		{
			"with brand filter",
			NewSelectParameters().WithBrand("Chanel"),
			func() string {
				choosingQuery := strings.TrimSpace(queries.SelectPerfumesBaseInfo) + " WHERE pb.brand = $1 AND s.sex = 'unisex'"
				withClause := strings.Replace(queries.WithSelect, "%s", choosingQuery, 1)
				return withClause + queries.EnrichSelectedPerfumes
			}(),
		},
		{
			"with sex filter",
			NewSelectParameters().WithSex("female"),
			func() string {
				choosingQuery := strings.TrimSpace(queries.SelectPerfumesBaseInfo) + " WHERE (s.sex = $1 OR s.sex = 'unisex')"
				withClause := strings.Replace(queries.WithSelect, "%s", choosingQuery, 1)
				return withClause + queries.EnrichSelectedPerfumes
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.GetQuery()
			if got != tt.want {
				t.Errorf("GetQuery() = %q, want %q", got, tt.want)
			}
		})
	}
}
