package connect

import (
	"Darkyfun/UrlShortener/internal/logging"
	"context"
	"fmt"
	"time"
)

const timeInterval = time.Second * 2

type Cacher interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (string, error)
	Ping(ctx context.Context) error
	Close() error
}

type Storage interface {
	GetAlias(ctx context.Context, orig string) string
	GetOriginal(ctx context.Context, alias string) (string, error)
	Set(ctx context.Context, alias string, orig string) error
	Ping(ctx context.Context) error
	Close()
}

func PingStorage(s Storage, logger *logging.EventLogger) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), timeInterval)

		err := s.Ping(ctx)
		if err != nil {
			fmt.Printf("ping database timeout with %v interval: %v\n", timeInterval, fmt.Errorf("connection closed or connection error"))
			logger.Log("warn", "ping database timeout: "+fmt.Sprint(timeInterval))
		}
		time.Sleep(timeInterval)
		cancel()
	}
}

func PingCache(c Cacher, logger *logging.EventLogger) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), timeInterval)

		err := c.Ping(ctx)
		if err != nil {
			fmt.Printf("ping cache timeout with %v interval: %v\n", timeInterval, fmt.Errorf("connection closed or connection error"))
			logger.Log("warn", "ping cache timeout: "+fmt.Sprint(timeInterval))
		}
		time.Sleep(timeInterval)
		cancel()
	}
}
