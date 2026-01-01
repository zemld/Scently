package perfume

import "github.com/zemld/Scently/models"

type Suggestions struct {
	Perfumes []models.Ranked `json:"suggested"`
}
