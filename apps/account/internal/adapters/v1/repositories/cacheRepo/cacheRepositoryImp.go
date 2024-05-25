package cacheRepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kurdilesmana/go-saving-api/apps/account/internal/core/ports/cachePort"
	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/kurdilesmana/go-saving-api/pkg/mid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	expiredMinnute = 10
	timeOut        = time.Second * 6
)

type CacheRepository struct {
	Cache          *redis.Client
	KeyTransaction string
	timeout        time.Duration
	log            *logging.Logger
}

func NewCacheRepo(cache *redis.Client, keyTransaction string, timeout int, log *logging.Logger) cachePort.ICacheRepository {
	return &CacheRepository{
		Cache:          cache,
		KeyTransaction: keyTransaction,
		timeout:        time.Duration(timeout) * time.Second,
		log:            log,
	}
}

func (r *CacheRepository) SetValue(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	ctxWT, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	return r.Cache.Set(ctxWT, key, value, ttl).Err()
}

func (r *CacheRepository) GetValue(ctx context.Context, key string) (string, error) {
	return r.Cache.Get(ctx, key).Result()
}

func (r *CacheRepository) DeleteValue(ctx context.Context, key string) error {
	ctxWT, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	return r.Cache.Del(ctxWT, key).Err()
}

func (r *CacheRepository) PublishStreamEvent(ctx context.Context, stream string, data map[string]interface{}) (err error) {
	ctxWT, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	r.log.Info(logrus.Fields{}, data, "Publishing event to Redis")
	jsonData, err := json.Marshal(data)
	if err != nil {
		r.log.Warn(logrus.Fields{}, data, err.Error())
		err = fmt.Errorf("failed to marshal publish data")
		return
	}

	requestID := mid.GetIDx(ctxWT)
	err = r.Cache.XAdd(ctxWT, &redis.XAddArgs{
		Stream: "transactions",
		MaxLen: 0,
		ID:     "",
		Values: map[string]interface{}{
			"whatHappened": string(stream + " received"),
			"eventID":      requestID,
			"eventData":    jsonData,
		},
	}).Err()

	return
}
