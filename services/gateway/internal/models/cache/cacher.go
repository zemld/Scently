package cache

import "context"

type Saver interface {
	Save(ctx context.Context, key string, value any) error
}

type Loader interface {
	Load(ctx context.Context, key string) (any, error)
}
