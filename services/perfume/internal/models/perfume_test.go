package models

import (
	"reflect"
	"testing"
)

func TestPerfume_Canonize(t *testing.T) {
	tests := []struct {
		name     string
		perfume  Perfume
		expected CanonizedPerfume
	}{
		{
			"simple perfume",
			Perfume{Brand: "Chanel", Name: "No.5"},
			CanonizedPerfume{Brand: "chanel", Name: "no5"},
		},
		{
			"perfume with spaces",
			Perfume{Brand: "Tom Ford", Name: "Oud Wood"},
			CanonizedPerfume{Brand: "tomford", Name: "oudwood"},
		},
		{
			"perfume with special chars",
			Perfume{Brand: "Dior-Sauvage!", Name: "Eau de Toilette"},
			CanonizedPerfume{Brand: "diorsauvage", Name: "eaudetoilette"},
		},
		{
			"perfume with numbers",
			Perfume{Brand: "Brand123", Name: "Name456"},
			CanonizedPerfume{Brand: "brand123", Name: "name456"},
		},
		{
			"uppercase perfume",
			Perfume{Brand: "CHANEL", Name: "NO.5"},
			CanonizedPerfume{Brand: "chanel", Name: "no5"},
		},
		{
			"mixed case perfume",
			Perfume{Brand: "ChAnEl", Name: "No.5"},
			CanonizedPerfume{Brand: "chanel", Name: "no5"},
		},
		{
			"empty strings",
			Perfume{Brand: "", Name: ""},
			CanonizedPerfume{Brand: "", Name: ""},
		},
		{
			"only special chars",
			Perfume{Brand: "!@#$", Name: "&*()"},
			CanonizedPerfume{Brand: "", Name: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.perfume.Canonize()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Perfume.Canonize() = %+v, want %+v", got, tt.expected)
			}
		})
	}
}
