package application

import (
	"log"

	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/matching"
	"github.com/zemld/PerfumeRecommendationSystem/algorithms/internal/domain/models"
)

type Tests struct {
	Brand string
	Name  string
	Sex   string
}

var tests = []Tests{
	{
		Brand: "DIOR",
		Name:  "J'adore",
		Sex:   "female", // женский из летуаль
	},
	{
		Brand: "Dolce & Gabbana",
		Name:  "Devotion Eau De Parfum Intense",
		Sex:   "female", // без семейств, по одной ноте
	},
	{
		Brand: "Dolce & Gabbana",
		Name:  "Dolce Lily",
		Sex:   "female", // без семейств, нот больше
	},
	{
		Brand: "1907",
		Name:  "Cedar Blue",
		Sex:   "male", // хороший
	},
	{
		Brand: "Abercrombie & Fitch",
		Name:  "Authentic Self Man",
		Sex:   "male", // без семейств
	},
	{
		Brand: "Adidas",
		Name:  "UEFA Champions League Champions Intense",
		Sex:   "male", // немного нот
	},
	{
		Brand: "100BON",
		Name:  "Ambre Sensuel",
		Sex:   "unisex", // по одной ноте
	},
	{
		Brand: "Versace",
		Name:  "Versus Uomo",
		Sex:   "male", // много нот
	},
	{
		Brand: "Aramis",
		Name:  "900",
		Sex:   "male", // мужской без семейств
	},
}

func RunTests(matcher matching.Matcher, allPerfumes []models.Perfume) ([]models.Perfume, [][]models.Ranked) {
	favs := make([]models.Perfume, 0, len(tests))
	results := make([][]models.Ranked, 0, len(tests))

	for _, test := range tests {
		favourite, ok := SelectConcretePerfume(test.Brand, test.Name, test.Sex, allPerfumes)
		if !ok {
			log.Printf("Can't get fav perfume %+v", test)
			continue
		}
		favs = append(favs, favourite)
		suggestions := matching.Find(*matching.NewMatchData(
			matcher,
			favourite,
			SelectPerfumesBySex(favourite.Sex, allPerfumes),
			4,
			6,
		))
		results = append(results, suggestions)
	}
	return favs, results
}
