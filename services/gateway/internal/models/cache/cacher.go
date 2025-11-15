package cache

import "context"

type Saver interface {
	Save(ctx context.Context, key string, value []byte) error
}

type Loader interface {
	Load(ctx context.Context, key string) ([]byte, error)
}
