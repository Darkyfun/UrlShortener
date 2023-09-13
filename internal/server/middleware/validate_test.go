package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		statusCode int
		body       string
		exp        string
	}{
		{name: "correct request", method: http.MethodPost, statusCode: http.StatusOK, body: `{"url":"https://www.google.com"}`, exp: "OK"},
		{name: "invalid body request", method: http.MethodPost, statusCode: http.StatusBadRequest, body: "invalid body", exp: "invalid request"},
		{name: "invalid json key request", method: http.MethodPost, statusCode: http.StatusBadRequest, body: `{"code":"https://www.google.com"}`, exp: "invalid request"},
		{name: "invalid url", method: http.MethodPost, statusCode: http.StatusBadRequest, body: `{"code":"https:/www.google.com"}`, exp: "invalid request"},
	}

	router := gin.New()
	router.Use(Validate())
	router.POST("/receive", func(context *gin.Context) {
		context.String(http.StatusOK, "OK")
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(tt.method, "/receive", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			assert.Nil(t, err)
			assert.Equal(t, tt.statusCode, w.Result().StatusCode)
			assert.Equal(t, tt.exp, w.Body.String())
		})
	}
}
