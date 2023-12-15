package middleware

import (
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var ErrInvalidRequest = errors.New("invalid request")
var ErrInvalidUrl = errors.New("invalid url")

// request - структура, предназначенная для парсинга JSON входящего запроса.
type request struct {
	Url string `json:"url"`
}

// Validate валидирует содержимное входящего http-запроса.
func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var r request
		err := c.BindJSON(&r)
		if err != nil {
			c.Set("status code", http.StatusBadRequest)
			c.String(http.StatusBadRequest, "invalid request")
			c.Abort()
			return
		}

		if r.Url == "" {
			c.String(http.StatusBadRequest, "%s", ErrInvalidRequest)
			c.Set("status code", http.StatusBadRequest)
			c.Abort()
			return
		}

		if ok := govalidator.IsURL(r.Url); ok == false {
			c.String(http.StatusBadRequest, "%s", ErrInvalidUrl)
			c.Set("status code", http.StatusBadRequest)
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		c.Request = c.Request.WithContext(context.WithValue(ctx, "IncomeUrl", r.Url))
		c.Next()
	}
}
