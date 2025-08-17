package internal

import (
	"context"
	"time"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/db/config"
)

func CreateContext(config *config.Config) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithValue(context.Background(), "dbConfig", config), 5*time.Second)
}
