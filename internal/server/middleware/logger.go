package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type LogHandler struct {
	ZapLog *zap.SugaredLogger
}

func NewLogHandler(file *os.File) *LogHandler {
	conf := zapcore.EncoderConfig{
		MessageKey:     "message",
		TimeKey:        "time",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("Mon, 02 Jan 2006 15:04:05.999 -0700"),
		EncodeDuration: zapcore.MillisDurationEncoder,
	}

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(conf), zapcore.Lock(file), zapcore.DebugLevel)
	logger := zap.New(core).Sugar()
	return &LogHandler{ZapLog: logger}
}

func (l *LogHandler) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			code, _ := c.Get("status code")
			l.ZapLog.Debugf("%v %v %v %v %v", c.Request.Method, c.Request.URL, c.Request.ContentLength, code, time.Since(start))
		}()
		c.Next()
	}
}
