package cachePort

import (
	"context"
	"time"
)

type ICacheRepository interface {
	SetValue(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error)
	GetValue(ctx context.Context, key string) (string, error)
	DeleteValue(ctx context.Context, key string) error
	PublishStreamEvent(ctx context.Context, stream string, data map[string]interface{}) error
}
