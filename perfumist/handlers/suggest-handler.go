package handlers

import (
	"net/http"

	"github.com/zemld/PerfumeRecommendationSystem/perfumist/internal"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/perplexity"
	"github.com/zemld/PerfumeRecommendationSystem/perfumist/util"
)

// 200 я получаю список духов
// 204 я не получаю ничего
// 400 я не передал бренд или название
// 500 какой-то трабл случился
// /v1/suggest/perfume?brand=string&name=string
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
	// TODO: обновить респонс (поле Success) и тут же может быть 500 ошибка

	filtered := filterSuggests(input.Brand, input.Name, suggested)
	var result []rankedPerfume
	for _, suggestedPerfume := range filtered {
		p := util.NewGetParameters().WithBrand(suggestedPerfume.Brand).WithName(suggestedPerfume.Name)
		suggestedPerfumesWithProps, ok := internal.GetPerfumes(*p)
		if !ok {
			util.WriteResponse(w, suggestResponse, http.StatusInternalServerError)
			return
		}
		if suggestedPerfumesWithProps == nil {
			continue
		}
		glued := util.GluePerfumes(suggestedPerfumesWithProps)
		result = append(result, rankedPerfume{Perfume: glued[0], Rank: suggestedPerfume.Rank})
	}
	if len(result) == 0 {
		util.WriteResponse(w, suggestResponse, http.StatusNoContent)
		return
	}
	suggestResponse.Suggested = result
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

// func GetSuggestedPerfumesProperties()
