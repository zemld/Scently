package matching

import "github.com/zemld/Scently/models"

type Matcher interface {
	Find(favourite models.Perfume, all []models.Perfume, matchesCount int) []models.Ranked
}
