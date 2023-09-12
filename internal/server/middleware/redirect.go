package middleware

import (
	"Darkyfun/UrlShortener/internal/logging"
	"Darkyfun/UrlShortener/internal/server/connect"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type aliasRequest struct {
	Alias string `uri:"alias" binding:"required"`
}

// Redirect handler for finding the original address by alias and redirecting
func Redirect(cache connect.Cacher, store connect.Storage, logger *logging.EventLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var q aliasRequest

		err := c.ShouldBindUri(&q)
		if err != nil {
			c.Set("status code", http.StatusBadRequest)
			c.String(http.StatusBadRequest, "invalid request")
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		orig, err := cache.Get(ctx, q.Alias)
		if err != nil && err.Error() == "cache miss" {

			orig, err = store.GetOriginal(ctx, q.Alias)
			if orig != "" && err == nil {

				if err = cache.Set(ctx, q.Alias, orig); err != nil {
					logger.Log("error", "reading and writing to cache failed")
					logger.Log("error", err.Error())
				}

				c.Set("status code", http.StatusTemporaryRedirect)
				c.Redirect(http.StatusTemporaryRedirect, orig)
			}
			//err = cache.Set(ctx, q.Alias, orig)
			if orig == "" {
				c.Set("status code", http.StatusBadRequest)
				c.String(http.StatusBadRequest, "Not found")
			}
		}

		c.Set("status code", http.StatusTemporaryRedirect)
		c.Redirect(http.StatusTemporaryRedirect, orig)
	}
}
