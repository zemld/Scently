package grpc_handlers

import (
	protoModels "github.com/zemld/Scently/generated/proto/perfume-hub/models"
	"github.com/zemld/Scently/perfume-hub/internal/models"
)

func convertPerfumeToProto(perfume models.Perfume) *protoModels.Perfume {
	return &protoModels.Perfume{
		Brand:      perfume.Brand,
		Name:       perfume.Name,
		Sex:        perfume.Sex,
		ImageUrl:   &perfume.ImageUrl,
		Properties: convertPropertiesToProto(perfume.Properties),
		Shops:      convertShopInfoToProto(perfume.Shops),
	}
}

func convertPropertiesToProto(properties models.PerfumeProperties) *protoModels.Perfume_Properties {
	return &protoModels.Perfume_Properties{
		PerfumeType: properties.Type,
		Family:      properties.Family,
		UpperNotes:  properties.UpperNotes,
		CoreNotes:   properties.CoreNotes,
		BaseNotes:   properties.BaseNotes,
	}
}

func convertShopInfoToProto(shops []models.ShopInfo) []*protoModels.Perfume_ShopInfo {
	responseShops := make([]*protoModels.Perfume_ShopInfo, len(shops))
	for i := range shops {
		responseShops[i] = &protoModels.Perfume_ShopInfo{
			ShopName: shops[i].ShopName,
			Domain:   shops[i].Domain,
			Variants: convertVariantsToProto(shops[i].Variants),
		}
	}
	return responseShops
}

func convertVariantsToProto(variants []models.PerfumeVariant) []*protoModels.Perfume_ShopInfo_Variant {
	responseVariants := make([]*protoModels.Perfume_ShopInfo_Variant, len(variants))
	for i := range variants {
		responseVariants[i] = &protoModels.Perfume_ShopInfo_Variant{
			Volume: int32(variants[i].Volume),
			Link:   variants[i].Link,
			Price:  int32(variants[i].Price),
		}
	}
	return responseVariants
}

func convertPerfumeToModel(perfume *protoModels.Perfume) models.Perfume {
	if perfume == nil {
		return models.Perfume{}
	}
	return models.Perfume{
		Brand:      perfume.Brand,
		Name:       perfume.Name,
		Sex:        perfume.Sex,
		ImageUrl:   tryConvertPointer(perfume.ImageUrl),
		Properties: convertPropertiesToModel(perfume.Properties),
		Shops:      convertShopInfoToModel(perfume.Shops),
	}
}

func convertPropertiesToModel(properties *protoModels.Perfume_Properties) models.PerfumeProperties {
	if properties == nil {
		return models.PerfumeProperties{}
	}
	return models.PerfumeProperties{
		Type:       properties.PerfumeType,
		Family:     properties.Family,
		UpperNotes: properties.UpperNotes,
		CoreNotes:  properties.CoreNotes,
		BaseNotes:  properties.BaseNotes,
	}
}

func convertShopInfoToModel(shops []*protoModels.Perfume_ShopInfo) []models.ShopInfo {
	modelShops := make([]models.ShopInfo, len(shops))
	for i := range shops {
		if shops[i] == nil {
			continue
		}
		modelShops[i] = models.ShopInfo{
			ShopName: shops[i].ShopName,
			Domain:   shops[i].Domain,
			ImageUrl: tryConvertPointer(shops[i].ImageUrl),
			Variants: convertVariantsToModel(shops[i].Variants),
		}
	}
	return modelShops
}

func convertVariantsToModel(variants []*protoModels.Perfume_ShopInfo_Variant) []models.PerfumeVariant {
	modelVariants := make([]models.PerfumeVariant, len(variants))
	for i := range variants {
		if variants[i] == nil {
			continue
		}
		modelVariants[i] = models.PerfumeVariant{
			Volume: int(variants[i].Volume),
			Link:   variants[i].Link,
			Price:  int(variants[i].Price),
		}
	}
	return modelVariants
}

func tryConvertPointer(value *string) string {
	if value != nil {
		return *value
	}
	return ""
}
