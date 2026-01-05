package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/zemld/config-manager/pkg/cm"
	"github.com/zemld/config-manager/pkg/cm/rcm"
)

var (
	managerOnce sync.Once
	manager     cm.ConfigManager
)

func Manager() cm.ConfigManager {
	managerOnce.Do(func() {
		manager = rcm.NewRedisConfigManager(
			"perfumist",
			&redis.Options{
				Addr:     fmt.Sprintf("%s:%s", os.Getenv("CONFIG_STORAGE_HOST"), os.Getenv("CONFIG_STORAGE_PORT")),
				Password: os.Getenv("CONFIG_STORAGE_PASSWORD"),
			},
		)
	})
	return manager
}
