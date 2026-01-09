package models

import (
	"reflect"
	"strings"
	"testing"

	queries "github.com/zemld/Scently/perfume-hub/internal/db/query"
)

func TestCanonize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple lowercase", "chanel", "chanel"},
		{"uppercase to lowercase", "CHANEL", "chanel"},
		{"mixed case", "ChAnEl", "chanel"},
		{"with spaces", "Chanel No.5", "chanelno5"},
		{"with numbers", "No.5", "no5"},
		{"with special chars", "Dior-Sauvage!", "diorsauvage"},
		{"with punctuation", "Tom Ford: Oud Wood", "tomfordoudwood"},
		{"empty string", "", ""},
		{"only special chars", "!@#$%", ""},
		{"with unicode letters", "Café", "café"},
		{"mixed with spaces and numbers", "Brand 123 Name", "brand123name"},
		{"only numbers", "123", "123"},
		{"numbers and letters", "Brand123Name", "brand123name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := canonize(tt.input)
			if got != tt.expected {
				t.Errorf("canonize(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

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
			baseQuery + " WHERE sex = 'unisex' AND page_number = $1",
		},
		{
			"brand only",
			NewSelectParameters().WithBrand("Chanel"),
			baseQuery + " WHERE sex = 'unisex' AND canonized_brand = $1",
		},
		{
			"name only",
			NewSelectParameters().WithName("No.5"),
			baseQuery + " WHERE sex = 'unisex' AND canonized_name = $1",
		},
		{
			"sex unisex",
			NewSelectParameters().WithSex("unisex"),
			baseQuery + " WHERE sex = 'unisex' AND page_number = $1",
		},
		{
			"sex female",
			NewSelectParameters().WithSex("female"),
			baseQuery + " WHERE (sex = 'unisex' OR sex = $1) AND page_number = $2",
		},
		{
			"sex male",
			NewSelectParameters().WithSex("male"),
			baseQuery + " WHERE (sex = 'unisex' OR sex = $1) AND page_number = $2",
		},
		{
			"brand and name",
			NewSelectParameters().WithBrand("Dior").WithName("Sauvage"),
			baseQuery + " WHERE sex = 'unisex' AND canonized_brand = $1 AND canonized_name = $2",
		},
		{
			"brand and sex female",
			NewSelectParameters().WithBrand("Chanel").WithSex("female"),
			baseQuery + " WHERE (sex = 'unisex' OR sex = $1) AND canonized_brand = $2",
		},
		{
			"brand and sex unisex",
			NewSelectParameters().WithBrand("Chanel").WithSex("unisex"),
			baseQuery + " WHERE sex = 'unisex' AND canonized_brand = $1",
		},
		{
			"name and sex unisex",
			NewSelectParameters().WithName("No.5").WithSex("unisex"),
			baseQuery + " WHERE sex = 'unisex' AND canonized_name = $1",
		},
		{
			"name and sex female",
			NewSelectParameters().WithName("No.5").WithSex("female"),
			baseQuery + " WHERE (sex = 'unisex' OR sex = $1) AND canonized_name = $2",
		},
		{
			"brand, name and sex male",
			NewSelectParameters().WithBrand("Dior").WithName("Sauvage").WithSex("male"),
			baseQuery + " WHERE (sex = 'unisex' OR sex = $1) AND canonized_brand = $2 AND canonized_name = $3",
		},
		{
			"brand, name and sex unisex",
			NewSelectParameters().WithBrand("Dior").WithName("Sauvage").WithSex("unisex"),
			baseQuery + " WHERE sex = 'unisex' AND canonized_brand = $1 AND canonized_name = $2",
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
			[]any{0},
		},
		{
			"brand only",
			NewSelectParameters().WithBrand("Chanel"),
			[]any{"chanel", 0},
		},
		{
			"name only",
			NewSelectParameters().WithName("No.5"),
			[]any{"no5", 0},
		},
		{
			"sex unisex - no args",
			NewSelectParameters().WithSex("unisex"),
			[]any{0},
		},
		{
			"sex female",
			NewSelectParameters().WithSex("female"),
			[]any{"female", 0},
		},
		{
			"sex male",
			NewSelectParameters().WithSex("male"),
			[]any{"male", 0},
		},
		{
			"brand and name",
			NewSelectParameters().WithBrand("Dior").WithName("Sauvage"),
			[]any{"dior", "sauvage"},
		},
		{
			"brand and sex female",
			NewSelectParameters().WithBrand("Chanel").WithSex("female"),
			[]any{"female", "chanel"},
		},
		{
			"brand and sex unisex - no sex arg",
			NewSelectParameters().WithBrand("Chanel").WithSex("unisex"),
			[]any{"chanel", 0},
		},
		{
			"name and sex female",
			NewSelectParameters().WithName("No.5").WithSex("female"),
			[]any{"female", "no5"},
		},
		{
			"brand, name and sex male",
			NewSelectParameters().WithBrand("Dior").WithName("Sauvage").WithSex("male"),
			[]any{"male", "dior", "sauvage"},
		},
		{
			"brand, name and sex unisex - no sex arg",
			NewSelectParameters().WithBrand("Dior").WithName("Sauvage").WithSex("unisex"),
			[]any{"dior", "sauvage"},
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
				choosingQuery := strings.TrimSpace(queries.SelectPerfumesBaseInfo) + " WHERE sex = 'unisex' AND page_number = $1"
				withClause := strings.Replace(queries.WithSelect, "%s", choosingQuery, 1)
				return withClause + queries.EnrichSelectedPerfumes
			}(),
		},
		{
			"with brand filter",
			NewSelectParameters().WithBrand("Chanel"),
			func() string {
				choosingQuery := strings.TrimSpace(queries.SelectPerfumesBaseInfo) + " WHERE sex = 'unisex' AND canonized_brand = $1"
				withClause := strings.Replace(queries.WithSelect, "%s", choosingQuery, 1)
				return withClause + queries.EnrichSelectedPerfumes
			}(),
		},
		{
			"with sex filter",
			NewSelectParameters().WithSex("female"),
			func() string {
				choosingQuery := strings.TrimSpace(queries.SelectPerfumesBaseInfo) + " WHERE (sex = 'unisex' OR sex = $1) AND page_number = $2"
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
