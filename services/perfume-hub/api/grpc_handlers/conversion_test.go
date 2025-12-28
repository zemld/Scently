package grpc_handlers

import (
	"reflect"
	"testing"

	protoModels "github.com/zemld/Scently/generated/proto/perfume-hub/models"
	"github.com/zemld/Scently/perfume-hub/internal/models"
)

func TestConvertPerfumeToProto(t *testing.T) {
	imageUrl := "http://example.com/image.jpg"
	perfume := models.Perfume{
		Brand:    "Chanel",
		Name:     "No.5",
		Sex:      "female",
		ImageUrl: imageUrl,
		Properties: models.PerfumeProperties{
			Type:       "Eau de Parfum",
			Family:     []string{"Floral"},
			UpperNotes: []string{"Bergamot", "Lemon"},
			CoreNotes:  []string{"Rose"},
			BaseNotes:  []string{"Musk"},
		},
		Shops: []models.ShopInfo{
			{
				ShopName: "Gold Apple",
				Domain:   "goldapple.ru",
				Variants: []models.PerfumeVariant{
					{Volume: 100, Price: 5000, Link: "http://example.com/link"},
				},
			},
		},
	}

	result := convertPerfumeToProto(perfume)

	if result.Brand != perfume.Brand {
		t.Errorf("convertPerfumeToProto() Brand = %q, want %q", result.Brand, perfume.Brand)
	}
	if result.Name != perfume.Name {
		t.Errorf("convertPerfumeToProto() Name = %q, want %q", result.Name, perfume.Name)
	}
	if result.Sex != perfume.Sex {
		t.Errorf("convertPerfumeToProto() Sex = %q, want %q", result.Sex, perfume.Sex)
	}
	if result.ImageUrl == nil || *result.ImageUrl != imageUrl {
		t.Errorf("convertPerfumeToProto() ImageUrl = %v, want %q", result.ImageUrl, imageUrl)
	}
	if result.Properties == nil {
		t.Fatal("convertPerfumeToProto() Properties is nil")
	}
	if len(result.Shops) != 1 {
		t.Errorf("convertPerfumeToProto() Shops len = %d, want 1", len(result.Shops))
	}
}

func TestConvertPropertiesToProto(t *testing.T) {
	properties := models.PerfumeProperties{
		Type:       "Eau de Toilette",
		Family:     []string{"Woody", "Spicy"},
		UpperNotes: []string{"Bergamot"},
		CoreNotes:  []string{"Rose", "Jasmine"},
		BaseNotes:  []string{"Musk", "Sandalwood"},
	}

	result := convertPropertiesToProto(properties)

	if result.PerfumeType != properties.Type {
		t.Errorf("convertPropertiesToProto() PerfumeType = %q, want %q", result.PerfumeType, properties.Type)
	}
	if !reflect.DeepEqual(result.Family, properties.Family) {
		t.Errorf("convertPropertiesToProto() Family = %v, want %v", result.Family, properties.Family)
	}
	if !reflect.DeepEqual(result.UpperNotes, properties.UpperNotes) {
		t.Errorf("convertPropertiesToProto() UpperNotes = %v, want %v", result.UpperNotes, properties.UpperNotes)
	}
	if !reflect.DeepEqual(result.CoreNotes, properties.CoreNotes) {
		t.Errorf("convertPropertiesToProto() CoreNotes = %v, want %v", result.CoreNotes, properties.CoreNotes)
	}
	if !reflect.DeepEqual(result.BaseNotes, properties.BaseNotes) {
		t.Errorf("convertPropertiesToProto() BaseNotes = %v, want %v", result.BaseNotes, properties.BaseNotes)
	}
}

func TestConvertShopInfoToProto(t *testing.T) {
	shops := []models.ShopInfo{
		{
			ShopName: "Gold Apple",
			Domain:   "goldapple.ru",
			Variants: []models.PerfumeVariant{
				{Volume: 100, Price: 5000, Link: "http://example.com/link1"},
				{Volume: 50, Price: 3000, Link: "http://example.com/link2"},
			},
		},
		{
			ShopName: "Randewoo",
			Domain:   "randewoo.ru",
			Variants: []models.PerfumeVariant{
				{Volume: 100, Price: 4500, Link: "http://example.com/link3"},
			},
		},
	}

	result := convertShopInfoToProto(shops)

	if len(result) != len(shops) {
		t.Fatalf("convertShopInfoToProto() len = %d, want %d", len(result), len(shops))
	}

	for i, shop := range shops {
		if result[i].ShopName != shop.ShopName {
			t.Errorf("convertShopInfoToProto()[%d] ShopName = %q, want %q", i, result[i].ShopName, shop.ShopName)
		}
		if result[i].Domain != shop.Domain {
			t.Errorf("convertShopInfoToProto()[%d] Domain = %q, want %q", i, result[i].Domain, shop.Domain)
		}
		if len(result[i].Variants) != len(shop.Variants) {
			t.Errorf("convertShopInfoToProto()[%d] Variants len = %d, want %d", i, len(result[i].Variants), len(shop.Variants))
		}
	}
}

func TestConvertVariantsToProto(t *testing.T) {
	variants := []models.PerfumeVariant{
		{Volume: 100, Price: 5000, Link: "http://example.com/link1"},
		{Volume: 50, Price: 3000, Link: "http://example.com/link2"},
		{Volume: 30, Price: 2000, Link: "http://example.com/link3"},
	}

	result := convertVariantsToProto(variants)

	if len(result) != len(variants) {
		t.Fatalf("convertVariantsToProto() len = %d, want %d", len(result), len(variants))
	}

	for i, variant := range variants {
		if int(result[i].Volume) != variant.Volume {
			t.Errorf("convertVariantsToProto()[%d] Volume = %d, want %d", i, result[i].Volume, variant.Volume)
		}
		if result[i].Link != variant.Link {
			t.Errorf("convertVariantsToProto()[%d] Link = %q, want %q", i, result[i].Link, variant.Link)
		}
		if int(result[i].Price) != variant.Price {
			t.Errorf("convertVariantsToProto()[%d] Price = %d, want %d", i, result[i].Price, variant.Price)
		}
	}
}

func TestConvertPerfumeToModel(t *testing.T) {
	imageUrl := "http://example.com/image.jpg"
	protoPerfume := &protoModels.Perfume{
		Brand:    "Chanel",
		Name:     "No.5",
		Sex:      "female",
		ImageUrl: &imageUrl,
		Properties: &protoModels.Perfume_Properties{
			PerfumeType: "Eau de Parfum",
			Family:      []string{"Floral"},
			UpperNotes:  []string{"Bergamot"},
			CoreNotes:   []string{"Rose"},
			BaseNotes:   []string{"Musk"},
		},
		Shops: []*protoModels.Perfume_ShopInfo{
			{
				ShopName: "Gold Apple",
				Domain:   "goldapple.ru",
				Variants: []*protoModels.Perfume_ShopInfo_Variant{
					{Volume: 100, Price: 5000, Link: "http://example.com/link"},
				},
			},
		},
	}

	result := convertPerfumeToModel(protoPerfume)

	if result.Brand != protoPerfume.Brand {
		t.Errorf("convertPerfumeToModel() Brand = %q, want %q", result.Brand, protoPerfume.Brand)
	}
	if result.Name != protoPerfume.Name {
		t.Errorf("convertPerfumeToModel() Name = %q, want %q", result.Name, protoPerfume.Name)
	}
	if result.Sex != protoPerfume.Sex {
		t.Errorf("convertPerfumeToModel() Sex = %q, want %q", result.Sex, protoPerfume.Sex)
	}
	if result.ImageUrl != imageUrl {
		t.Errorf("convertPerfumeToModel() ImageUrl = %q, want %q", result.ImageUrl, imageUrl)
	}
	if len(result.Shops) != 1 {
		t.Errorf("convertPerfumeToModel() Shops len = %d, want 1", len(result.Shops))
	}
}

func TestConvertPerfumeToModel_NilPerfume(t *testing.T) {
	result := convertPerfumeToModel(nil)

	if result.Brand != "" {
		t.Errorf("convertPerfumeToModel(nil) Brand = %q, want empty", result.Brand)
	}
	if result.Name != "" {
		t.Errorf("convertPerfumeToModel(nil) Name = %q, want empty", result.Name)
	}
}

func TestConvertPropertiesToModel(t *testing.T) {
	protoProperties := &protoModels.Perfume_Properties{
		PerfumeType: "Eau de Toilette",
		Family:      []string{"Woody", "Spicy"},
		UpperNotes:  []string{"Bergamot"},
		CoreNotes:   []string{"Rose", "Jasmine"},
		BaseNotes:   []string{"Musk", "Sandalwood"},
	}

	result := convertPropertiesToModel(protoProperties)

	if result.Type != protoProperties.PerfumeType {
		t.Errorf("convertPropertiesToModel() Type = %q, want %q", result.Type, protoProperties.PerfumeType)
	}
	if !reflect.DeepEqual(result.Family, protoProperties.Family) {
		t.Errorf("convertPropertiesToModel() Family = %v, want %v", result.Family, protoProperties.Family)
	}
	if !reflect.DeepEqual(result.UpperNotes, protoProperties.UpperNotes) {
		t.Errorf("convertPropertiesToModel() UpperNotes = %v, want %v", result.UpperNotes, protoProperties.UpperNotes)
	}
	if !reflect.DeepEqual(result.CoreNotes, protoProperties.CoreNotes) {
		t.Errorf("convertPropertiesToModel() CoreNotes = %v, want %v", result.CoreNotes, protoProperties.CoreNotes)
	}
	if !reflect.DeepEqual(result.BaseNotes, protoProperties.BaseNotes) {
		t.Errorf("convertPropertiesToModel() BaseNotes = %v, want %v", result.BaseNotes, protoProperties.BaseNotes)
	}
}

func TestConvertPropertiesToModel_NilProperties(t *testing.T) {
	result := convertPropertiesToModel(nil)

	if result.Type != "" {
		t.Errorf("convertPropertiesToModel(nil) Type = %q, want empty", result.Type)
	}
	if result.Family != nil {
		t.Errorf("convertPropertiesToModel(nil) Family = %v, want nil", result.Family)
	}
}

func TestConvertShopInfoToModel(t *testing.T) {
	imageUrl1 := "http://example.com/image1.jpg"
	imageUrl2 := "http://example.com/image2.jpg"
	protoShops := []*protoModels.Perfume_ShopInfo{
		{
			ShopName: "Gold Apple",
			Domain:   "goldapple.ru",
			ImageUrl: &imageUrl1,
			Variants: []*protoModels.Perfume_ShopInfo_Variant{
				{Volume: 100, Price: 5000, Link: "http://example.com/link1"},
				{Volume: 50, Price: 3000, Link: "http://example.com/link2"},
			},
		},
		{
			ShopName: "Randewoo",
			Domain:   "randewoo.ru",
			ImageUrl: &imageUrl2,
			Variants: []*protoModels.Perfume_ShopInfo_Variant{
				{Volume: 100, Price: 4500, Link: "http://example.com/link3"},
			},
		},
	}

	result := convertShopInfoToModel(protoShops)

	if len(result) != len(protoShops) {
		t.Fatalf("convertShopInfoToModel() len = %d, want %d", len(result), len(protoShops))
	}

	for i, shop := range protoShops {
		if result[i].ShopName != shop.ShopName {
			t.Errorf("convertShopInfoToModel()[%d] ShopName = %q, want %q", i, result[i].ShopName, shop.ShopName)
		}
		if result[i].Domain != shop.Domain {
			t.Errorf("convertShopInfoToModel()[%d] Domain = %q, want %q", i, result[i].Domain, shop.Domain)
		}
		if result[i].ImageUrl != *shop.ImageUrl {
			t.Errorf("convertShopInfoToModel()[%d] ImageUrl = %q, want %q", i, result[i].ImageUrl, *shop.ImageUrl)
		}
		if len(result[i].Variants) != len(shop.Variants) {
			t.Errorf("convertShopInfoToModel()[%d] Variants len = %d, want %d", i, len(result[i].Variants), len(shop.Variants))
		}
	}
}

func TestConvertShopInfoToModel_NilShops(t *testing.T) {
	protoShops := []*protoModels.Perfume_ShopInfo{
		nil,
		{
			ShopName: "Gold Apple",
			Domain:   "goldapple.ru",
			Variants: []*protoModels.Perfume_ShopInfo_Variant{},
		},
	}

	result := convertShopInfoToModel(protoShops)

	// nil shops должны быть пропущены
	if len(result) != 2 {
		t.Fatalf("convertShopInfoToModel() len = %d, want 2", len(result))
	}
}

func TestConvertVariantsToModel(t *testing.T) {
	protoVariants := []*protoModels.Perfume_ShopInfo_Variant{
		{Volume: 100, Price: 5000, Link: "http://example.com/link1"},
		{Volume: 50, Price: 3000, Link: "http://example.com/link2"},
		{Volume: 30, Price: 2000, Link: "http://example.com/link3"},
	}

	result := convertVariantsToModel(protoVariants)

	if len(result) != len(protoVariants) {
		t.Fatalf("convertVariantsToModel() len = %d, want %d", len(result), len(protoVariants))
	}

	for i, variant := range protoVariants {
		if result[i].Volume != int(variant.Volume) {
			t.Errorf("convertVariantsToModel()[%d] Volume = %d, want %d", i, result[i].Volume, variant.Volume)
		}
		if result[i].Link != variant.Link {
			t.Errorf("convertVariantsToModel()[%d] Link = %q, want %q", i, result[i].Link, variant.Link)
		}
		if result[i].Price != int(variant.Price) {
			t.Errorf("convertVariantsToModel()[%d] Price = %d, want %d", i, result[i].Price, variant.Price)
		}
	}
}

func TestConvertVariantsToModel_NilVariants(t *testing.T) {
	protoVariants := []*protoModels.Perfume_ShopInfo_Variant{
		nil,
		{Volume: 100, Price: 5000, Link: "http://example.com/link"},
	}

	result := convertVariantsToModel(protoVariants)

	// nil variants должны быть пропущены
	if len(result) != 2 {
		t.Fatalf("convertVariantsToModel() len = %d, want 2", len(result))
	}
}

func TestTryConvertPointer(t *testing.T) {
	tests := []struct {
		name     string
		value    *string
		expected string
	}{
		{
			name:     "nil pointer",
			value:    nil,
			expected: "",
		},
		{
			name:     "non-nil pointer",
			value:    stringPtr("test"),
			expected: "test",
		},
		{
			name:     "empty string pointer",
			value:    stringPtr(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tryConvertPointer(tt.value)
			if result != tt.expected {
				t.Errorf("tryConvertPointer() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func TestConvertRoundTrip(t *testing.T) {
	// Тест на то, что конвертация туда-обратно сохраняет данные
	imageUrl := "http://example.com/image.jpg"
	original := models.Perfume{
		Brand:    "Chanel",
		Name:     "No.5",
		Sex:      "female",
		ImageUrl: imageUrl,
		Properties: models.PerfumeProperties{
			Type:       "Eau de Parfum",
			Family:     []string{"Floral"},
			UpperNotes: []string{"Bergamot", "Lemon"},
			CoreNotes:  []string{"Rose"},
			BaseNotes:  []string{"Musk"},
		},
		Shops: []models.ShopInfo{
			{
				ShopName: "Gold Apple",
				Domain:   "goldapple.ru",
				ImageUrl: imageUrl,
				Variants: []models.PerfumeVariant{
					{Volume: 100, Price: 5000, Link: "http://example.com/link"},
				},
			},
		},
	}

	proto := convertPerfumeToProto(original)
	converted := convertPerfumeToModel(proto)

	if converted.Brand != original.Brand {
		t.Errorf("RoundTrip Brand = %q, want %q", converted.Brand, original.Brand)
	}
	if converted.Name != original.Name {
		t.Errorf("RoundTrip Name = %q, want %q", converted.Name, original.Name)
	}
	if converted.Sex != original.Sex {
		t.Errorf("RoundTrip Sex = %q, want %q", converted.Sex, original.Sex)
	}
	if converted.ImageUrl != original.ImageUrl {
		t.Errorf("RoundTrip ImageUrl = %q, want %q", converted.ImageUrl, original.ImageUrl)
	}
	if converted.Properties.Type != original.Properties.Type {
		t.Errorf("RoundTrip Properties.Type = %q, want %q", converted.Properties.Type, original.Properties.Type)
	}
	if !reflect.DeepEqual(converted.Properties.Family, original.Properties.Family) {
		t.Errorf("RoundTrip Properties.Family = %v, want %v", converted.Properties.Family, original.Properties.Family)
	}
	if len(converted.Shops) != len(original.Shops) {
		t.Errorf("RoundTrip Shops len = %d, want %d", len(converted.Shops), len(original.Shops))
	}
}
