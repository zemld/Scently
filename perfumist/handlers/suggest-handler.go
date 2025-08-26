package handlers

import (
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/perplexity"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/util"
)

// @description Get suggests for perfumes. Accept brand and name and recommends 4+- perfumes which user probably will like.
// @tags Perfumes
// @summary Suggests some perfumes
// @produce json
// @param brand query string true "Brand of the perfume which you like"
// @param name query string true "Name of the perfume which you like"
// @success 200 {object} SuggestResponse "Suggested perfumes"
// @success 204 {object} SuggestResponse "No perfumes found for suggestion"
// @failure 400 {object} SuggestResponse "Incorrect parameters"
// @failure 500 {object} SuggestResponse
// @router /perfume [get]
func SuggestHandler(w http.ResponseWriter, r *http.Request) {
	var suggestResponse SuggestResponse
	input, ok := parseQuery(r, &suggestResponse)
	if !ok {
		util.WriteResponse(w, suggestResponse, http.StatusBadRequest)
		return
	}

	// perplexity
	// TODO: implement perplexity + check perfume existence
	_ = input

	var suggested []perplexity.RankedPerfume = nil
	suggested = append(suggested, perplexity.RankedPerfume{Brand: "Marc Jacobs", Name: "Daisy Love", Rank: 1})
	suggested = append(suggested, perplexity.RankedPerfume{Brand: "Jardin de Parfums", Name: "UNIQUE LOVE LETTER", Rank: 2})
	suggested = append(suggested, perplexity.RankedPerfume{Brand: "Chloé", Name: "Nomade", Rank: 3})
	// TODO: обновить респонс (поле Success) и тут же может быть 500 ошибка

	filtered := filterSuggests(input.Brand, input.Name, suggested)
	enriched, ok := enrichSuggestedPerfumes(filtered)
	if !ok {
		suggestResponse.Success = false
		util.WriteResponse(w, suggestResponse, http.StatusInternalServerError)
		return
	}
	if len(enriched) == 0 {
		util.WriteResponse(w, suggestResponse, http.StatusNoContent)
		return
	}
	suggestResponse.Suggested = enriched
	util.WriteResponse(w, suggestResponse, http.StatusOK)
}

func parseQuery(r *http.Request, suggestResponse *SuggestResponse) (util.GetParameters, bool) {
	brand := r.URL.Query().Get("brand")
	name := r.URL.Query().Get("name")
	suggestResponse.Input = inputPerfume{Brand: brand, Name: name, Ok: true}
	suggestResponse.Success = true
	if brand == "" || name == "" {
		suggestResponse.Input.Ok = false
		suggestResponse.Success = false
		return util.GetParameters{}, false
	}
	return *util.NewGetParameters().WithBrand(brand).WithName(name), true
}

func filterSuggests(inputBrand string, inputName string, suggests []perplexity.RankedPerfume) []perplexity.RankedPerfume {
	var filtered []perplexity.RankedPerfume
	currentRank := 1
	for _, s := range suggests {
		if s.Brand == inputBrand && s.Name == inputName {
			continue
		}
		f := s
		f.Rank = currentRank
		currentRank++
		filtered = append(filtered, f)
	}
	return filtered
}

func enrichSuggestedPerfumes(suggested []perplexity.RankedPerfume) ([]rankedPerfumeWithProps, bool) {
	var result []rankedPerfumeWithProps
	for _, suggestedPerfume := range suggested {
		p := util.NewGetParameters().WithBrand(suggestedPerfume.Brand).WithName(suggestedPerfume.Name)
		suggestedPerfumesWithProps, ok := internal.GetPerfumes(*p)
		if !ok {
			return nil, false
		}
		if suggestedPerfumesWithProps == nil {
			continue
		}
		glued := util.GluePerfumes(suggestedPerfumesWithProps)
		result = append(result, rankedPerfumeWithProps{Perfume: glued[0], Rank: suggestedPerfume.Rank})
	}
	return result, true
}
