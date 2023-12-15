package cache

import (
	"Darkyfun/UrlShortener/internal/logging"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"testing"
	"time"
)

const redisAddr = "localhost:6379"

const testKey = "test_key"
const testValue = "test_value"

func TestMain(m *testing.M) {
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	defer func() {
		err := rdb.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	rdb.Set(context.Background(), testKey, testValue, time.Hour)
	rdb.Del(context.Background(), "unknown key")

	m.Run()
}

func TestRapidDb_Set(t *testing.T) {
	rdb := NewCacheDb(Opts{Addr: "localhost:6379"}, logging.NewLogger("json", io.Discard))

	// happy logpath.
	err := rdb.Set(context.Background(), testKey, testValue)
	assert.Nil(t, err)

	// timeout.
	ctxExp, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(time.Millisecond)

	err = rdb.Set(ctxExp, testKey, testValue)
	assert.Equal(t, ErrFailed, err)

	// closed client.
	err = rdb.Close()
	if err != nil {
		t.Errorf("failed to close connection to cache during test: %v\n", err)
	}
	err = rdb.Set(context.Background(), testKey, testValue)
	assert.Equal(t, ErrClientClosed, err)
}

func TestRapidDb_Get(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		err   error
	}{
		{name: "happy logpath", key: testKey, value: testValue, err: nil},
		{name: "cache miss", key: "unknown key", value: "", err: ErrCacheMiss},
	}

	rdb := NewCacheDb(Opts{Addr: redisAddr}, logging.NewLogger("json", io.Discard))

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			res, err := rdb.Get(context.Background(), tt.key)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.value, res)
		})
	}

	// timeout.
	ctxExp, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(time.Millisecond)

	res, err := rdb.Get(ctxExp, testValue)
	assert.Equal(t, ErrFailed, err)
	assert.Equal(t, "", res)

	// closed client.
	err = rdb.Close()
	if err != nil {
		t.Errorf("failed to close connection to cache during test: %v\n", err)
	}
	res, err = rdb.Get(context.Background(), testKey)
	assert.Equal(t, ErrClientClosed, err)
	assert.Equal(t, "", res)
}
