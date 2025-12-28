package fetching

import (
	"context"
	"net/http"
	"sync"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/config"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/parameters"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal/models/perfume"
	ph "github.com/zemld/Scently/generated/proto/perfume-hub"
	protoModels "github.com/zemld/Scently/generated/proto/perfume-hub/models"
	"github.com/zemld/Scently/generated/proto/perfume-hub/requests"
)

type PerfumeHub struct {
	client ph.PerfumeStorageClient
}

func NewPerfumeHub(client ph.PerfumeStorageClient) *PerfumeHub {
	return &PerfumeHub{
		client: client,
	}
}

func (f *PerfumeHub) Fetch(
	ctx context.Context,
	params []parameters.RequestPerfume,
) ([]perfume.Perfume, bool) {
	perfumesChan := make(chan *requests.GetPerfumesResponse, len(params))

	wg := sync.WaitGroup{}
	wg.Add(len(params))

	ctx, cancel := context.WithTimeout(ctx, config.PerfumeHubFetcherTimeout)
	defer cancel()

	for _, param := range params {
		go f.getPerfumesAsync(ctx, param, perfumesChan, &wg)
	}

	go func() {
		wg.Wait()
		close(perfumesChan)
	}()

	all, AllStatus := f.fetchPerfumeResults(ctx, perfumesChan)
	if AllStatus == http.StatusForbidden || AllStatus >= http.StatusInternalServerError {
		return nil, false
	}
	if len(all) == 0 {
		return nil, false
	}
	return all, true
}

func (f *PerfumeHub) getPerfumesAsync(
	ctx context.Context,
	params parameters.RequestPerfume,
	results chan<- *requests.GetPerfumesResponse,
	wg *sync.WaitGroup,
) {
	if wg != nil {
		defer wg.Done()
	}
	response, err := f.client.GetPerfumes(ctx, &requests.GetPerfumesRequest{
		Brand: params.Brand,
		Name:  params.Name,
		Sex:   params.Sex,
	})
	if err != nil {
		results <- nil
		return
	}
	results <- response
}

func (f *PerfumeHub) fetchPerfumeResults(
	ctx context.Context,
	perfumesChan <-chan *requests.GetPerfumesResponse,
) ([]perfume.Perfume, int) {
	var perfumes []perfume.Perfume
	var status int

	for {
		select {
		case result, ok := <-perfumesChan:
			if !ok {
				return perfumes, status
			}
			if result == nil {
				continue
			}
			for _, perfume := range result.Perfumes {
				if perfume == nil {
					continue
				}
				perfumes = append(perfumes, convertPerfumeToModel(perfume))
			}
		case <-ctx.Done():
			return nil, http.StatusInternalServerError
		}
	}
}

func convertPerfumeToModel(p *protoModels.Perfume) perfume.Perfume {
	if p == nil {
		return perfume.Perfume{}
	}
	return perfume.Perfume{
		Brand:      p.Brand,
		Name:       p.Name,
		Sex:        p.Sex,
		ImageUrl:   tryConvertPointer(p.ImageUrl),
		Properties: convertPropertiesToModel(p.Properties),
		Shops:      convertShopInfoToModel(p.Shops),
	}
}

func convertPropertiesToModel(properties *protoModels.Perfume_Properties) perfume.Properties {
	if properties == nil {
		return perfume.Properties{}
	}
	return perfume.Properties{
		Type:       properties.PerfumeType,
		Family:     properties.Family,
		UpperNotes: properties.UpperNotes,
		CoreNotes:  properties.CoreNotes,
		BaseNotes:  properties.BaseNotes,
	}
}

func convertShopInfoToModel(shops []*protoModels.Perfume_ShopInfo) []perfume.ShopInfo {
	modelShops := make([]perfume.ShopInfo, len(shops))
	for i := range shops {
		if shops[i] == nil {
			continue
		}
		modelShops[i] = perfume.ShopInfo{
			ShopName: shops[i].ShopName,
			Domain:   shops[i].Domain,
			ImageUrl: tryConvertPointer(shops[i].ImageUrl),
			Variants: convertVariantsToModel(shops[i].Variants),
		}
	}
	return modelShops
}

func convertVariantsToModel(variants []*protoModels.Perfume_ShopInfo_Variant) []perfume.Variant {
	modelVariants := make([]perfume.Variant, len(variants))
	for i := range variants {
		if variants[i] == nil {
			continue
		}
		modelVariants[i] = perfume.Variant{
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
