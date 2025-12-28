package core

import (
	"context"
	"time"

	"github.com/zemld/Scently/perfume-hub/internal/db/config"
)

func CreateContext(c *config.Config) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithValue(context.Background(), config.ConfigKey, c), 5*time.Second)
}
