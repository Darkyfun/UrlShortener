package cache

import (
	"Darkyfun/UrlShortener/internal/logging"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrFailed = errors.New("operation failed")
var ErrClientClosed = errors.New("unable to set a record: client is closed")
var ErrCacheMiss = errors.New("cache miss")

type RapidDb struct {
	rdb    *redis.Client
	logger *logging.EventLogger
}

type Opts struct {
	Addr, User, Password string
	MaxRetries, PoolSize int
}

func NewCacheDb(options Opts, log *logging.EventLogger) *RapidDb {
	client := redis.NewClient(&redis.Options{
		Addr: options.Addr,
	})

	return &RapidDb{rdb: client, logger: log}
}

func (c *RapidDb) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *RapidDb) Close() error {
	return c.rdb.Close()
}

// Set sets key-value pair in the cache
func (c *RapidDb) Set(ctx context.Context, key string, value any) error {
	_, err := c.rdb.Set(ctx, key, value, time.Hour).Result()
	if err != nil && err.Error() == "redis: client is closed" {
		c.logger.Log("error", ErrClientClosed.Error())
		return ErrClientClosed
	} else if err != nil {
		c.logger.Log("error", ErrFailed.Error())
		return ErrFailed
	}
	return nil
}

// Get gets value by provided key
func (c *RapidDb) Get(ctx context.Context, key string) (string, error) {
	res := c.rdb.Get(ctx, key)
	url, err := res.Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return "", ErrCacheMiss
		} else {
			return "", ErrFailed
		}
	}
	return url, nil
}