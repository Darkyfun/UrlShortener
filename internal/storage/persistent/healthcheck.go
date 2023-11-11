package persistent

import (
	"Darkyfun/UrlShortener/internal/logging"
	"context"
	"fmt"
	"time"
)

const timeInterval = time.Second * 2

type Pinger interface {
	Ping(ctx context.Context) error
}

func PingStorage(s Pinger, log logging.Logger) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), timeInterval)

		err := s.Ping(ctx)
		if err != nil {
			fmt.Printf("ping database timeout with %v interval: %v\n", timeInterval, fmt.Errorf("connection closed or connection error"))
			log.Log("warn", "ping database timeout: "+fmt.Sprint(timeInterval))
		}
		time.Sleep(timeInterval)
		cancel()
	}
}
