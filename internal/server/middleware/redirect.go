package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Logger interface {
	Log(string, string)
}

// aliasRequest - это структура, предназначенная для парсинга URL-параметра входящего запроса.
type aliasRequest struct {
	Alias string `uri:"alias" binding:"required"`
}

// Redirect парсит входящий запрос с псевдонимом и перенаправляет клиент на оригинальный URL с кодом ответа 307.
// Сначала Redirect проверят кэш на наличие записи. Если данная запись есть, то осуществляется перенаправление.
// Если в кэше записи нет, то запрос на выборку отправляется в SQL-базу данных, после чего в кэш вносится данная пара значений и клиента перенаправляют на оригинальный URL.
func Redirect(cache Cacher, store Storager, logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var q aliasRequest
		err := c.ShouldBindUri(&q)
		if err != nil {
			c.Set("status code", http.StatusBadRequest)
			c.String(http.StatusBadRequest, "invalid request")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		orig, err := cache.Get(ctx, q.Alias)
		if err != nil && err.Error() == "cache miss" {

			orig, err = store.GetOriginal(ctx, q.Alias)
			if orig != "" && err == nil {

				if err = cache.Set(ctx, q.Alias, orig); err != nil {
					logger.Log("error", "reading and writing to cache failed: "+err.Error())
				}

				c.Set("status code", http.StatusTemporaryRedirect)
				c.Redirect(http.StatusTemporaryRedirect, orig)
				return
			}
			if orig == "" {
				c.Set("status code", http.StatusBadRequest)
				c.String(http.StatusBadRequest, "Not found")
				return
			}
		}

		c.Set("status code", http.StatusTemporaryRedirect)
		c.Redirect(http.StatusTemporaryRedirect, orig)
	}
}
