package middleware

import (
	"Darkyfun/UrlShortener/internal/lib"
	"Darkyfun/UrlShortener/internal/server/connect"
	"Darkyfun/UrlShortener/internal/storage/persistent"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Saver is handler who parses original url and then responses with alias url
func Saver(cache connect.Cacher, store connect.Storage, addr string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origUrl := c.Request.Context().Value("IncomeUrl").(string)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		alias := store.GetAlias(ctx, origUrl)
		if alias != "" {
			c.Set("status code", http.StatusOK)
			c.JSON(http.StatusOK, gin.H{
				"Short_url": "http://" + "localhost" + addr + "/redirect/" + alias,
			})
			return
		}

		alias = lib.GetRandomAlias(10)
		for {
			err := store.Set(ctx, alias, origUrl)
			if err == nil || errors.Is(err, persistent.ErrConnClosed) {
				break
			}
		}
		err := cache.Set(ctx, alias, origUrl)
		if err != nil {
			c.Set("status code", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		c.Set("status code", http.StatusOK)
		c.JSON(http.StatusOK, gin.H{
			"Short_url": "http://" + "localhost" + addr + "/redirect/" + alias,
		})
	}
}
