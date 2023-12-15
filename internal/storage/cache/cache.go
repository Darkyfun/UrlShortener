// Package cache представляет собой реализацию методов для работы с базой данных, являющейся кэшем.
package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var ErrFailed = errors.New("operation failed")
var ErrClientClosed = errors.New("unable to set a record: client is closed")
var ErrCacheMiss = errors.New("cache miss")

// RapidDb - это структура, реализующая запросы к базе данных, являющейся кэшем.
type RapidDb struct {
	rdb *redis.Client
	log Logger
}

// Opts - это опции, необходимые для подключения к кэшу.
type Opts struct {
	Addr, User, Password string
	MaxRetries, PoolSize int
}

// NewCacheDb возвращает переменную *RapidDb, готовую к работе с кэшем.
func NewCacheDb(options Opts, logg Logger) *RapidDb {
	client := redis.NewClient(&redis.Options{
		Addr: options.Addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	return &RapidDb{rdb: client, log: logg}
}

// Ping тестирует соединение с кэшем.
func (c *RapidDb) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Close закрывает соединение с кэшем.
func (c *RapidDb) Close() error {
	return c.rdb.Close()
}

// Set сохраняет в кэше запись, состояющую из псевдонима и оригинального URL.
func (c *RapidDb) Set(ctx context.Context, keyAlias string, valueOriginal any) error {
	_, err := c.rdb.Set(ctx, keyAlias, valueOriginal, time.Hour).Result()

	if err != nil && err.Error() == "redis: client is closed" {
		c.log.Log("error", ErrClientClosed.Error())
		return ErrClientClosed
	} else if err != nil {
		c.log.Log("error", ErrFailed.Error())
		return ErrFailed
	}

	return nil
}

// Get получает значения псевдонима из кэша по указанному псевдониму
func (c *RapidDb) Get(ctx context.Context, keyAlias string) (string, error) {
	res := c.rdb.Get(ctx, keyAlias)
	url, err := res.Result()

	if err != nil {
		switch {
		case err.Error() == "redis: nil":
			return "", ErrCacheMiss
		case err.Error() == "redis: client is closed":
			return "", ErrClientClosed
		default:
			return "", ErrFailed
		}
	}

	return url, nil
}
