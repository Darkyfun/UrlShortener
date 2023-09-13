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
	"testing"
)

const TestBase = "postgres://go:123@localhost:5432/service_test"
const redisAddr = "localhost:6379"

func TestRedirect(t *testing.T) {
	tests := []struct {
		name string
		url  string
		resp int
	}{
		{name: "valid request", url: "/redirect/googlealias", resp: http.StatusTemporaryRedirect},
		{name: "no alias", url: "/redirect", resp: http.StatusBadRequest},
		{name: "invalid alias", url: "/redirect/invalidalias", resp: http.StatusBadRequest},
	}

	logger := logging.NewLogger("json", io.Discard)
	db := persistent.NewDb(context.Background(), logger, TestBase)
	rdb := cache.NewCacheDb(cache.Opts{Addr: redisAddr}, logger)

	router := gin.New()
	router.GET("/redirect/:alias", Redirect(rdb, db, logger))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/redirect/googlealias", nil)

			router.ServeHTTP(w, r)
			assert.Equal(t, http.StatusTemporaryRedirect, w.Result().StatusCode)
		})
	}
}
