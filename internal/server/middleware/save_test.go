package middleware

import (
	"Darkyfun/UrlShortener/internal/logging"
	"Darkyfun/UrlShortener/internal/storage/cache"
	"Darkyfun/UrlShortener/internal/storage/persistent"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//const TestBase = "postgres://go:123@localhost:5432/service_test"
//const redisAddr = "localhost:6379"

func TestSaver(t *testing.T) {
	logger := logging.NewLogger("json", io.Discard)
	db := persistent.NewDb(context.Background(), logger, TestBase)
	rdb := cache.NewCacheDb(cache.Opts{Addr: redisAddr}, logger)

	// проверяем доступность кэша и базы данных.
	err := db.Ping(context.Background())
	assert.Nil(t, err)
	err = rdb.Ping(context.Background())
	assert.Nil(t, err)

	// настраиваем gin.
	router := gin.New()
	router.Use(Saver(rdb, &db, ":5050"))
	router.POST("/receive")

	// happy path.
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/receive", strings.NewReader("{\"url\":\"https://www.google.com\"}"))
	assert.Nil(t, err)
	r = r.WithContext(context.WithValue(context.Background(), "IncomeUrl", "https://www.google.com"))
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	// закрываем соединение с кэшем.
	db.Close()
	err = rdb.Close()

	// тест с закрытым соединение к кэшу.
	w = httptest.NewRecorder()
	r, err = http.NewRequest(http.MethodPost, "/receive", strings.NewReader("{\"url\":\"https://www.google.com\"}"))
	assert.Nil(t, err)
	r = r.WithContext(context.WithValue(context.Background(), "IncomeUrl", "https://www.google.com"))
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}
