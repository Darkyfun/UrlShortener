package persistent

import (
	"context"
	"fmt"
	"time"
)

const timeInterval = time.Second * 2

type Pinger interface {
	Ping(ctx context.Context) error
}

type Logger interface {
	Log(string, string)
}

// PingStorage циклично отправляет ping-запросы в базу данных для её мониторинга.
func PingStorage(s Pinger, log Logger) {
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
