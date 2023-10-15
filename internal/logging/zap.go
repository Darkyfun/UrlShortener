package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

type Logger interface {
	Log(string, string)
}

type EventLogger struct {
	logger *zap.Logger
}

func (c *EventLogger) Log(level string, msg string) {
	switch level {
	case "debug":
		c.logger.Log(zap.DebugLevel, msg)
	case "info":
		c.logger.Log(zap.InfoLevel, msg)
	case "warn":
		c.logger.Log(zap.WarnLevel, msg)
	case "error":
		c.logger.Log(zap.ErrorLevel, msg)
	case "panic":
		c.logger.Log(zap.PanicLevel, msg)
	case "fatal":
		c.logger.Log(zap.FatalLevel, msg)
	}
}

func NewLogger(outputType string, file io.Writer) *EventLogger {
	var encType zapcore.Encoder
	
	conf := zapcore.EncoderConfig{
		MessageKey:          "message",
		LevelKey:            "log_level",
		TimeKey:             "time",
		StacktraceKey:       "trace",
		SkipLineEnding:      false,
		LineEnding:          ";\n",
		EncodeLevel:         zapcore.CapitalLevelEncoder,
		EncodeTime:          zapcore.TimeEncoderOfLayout("Mon, 02 Jan 2006 15:04:05.999 -0700"),
		EncodeDuration:      zapcore.MillisDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		EncodeName:          zapcore.FullNameEncoder,
		NewReflectedEncoder: nil,
	}

	if outputType == "json" {
		encType = zapcore.NewJSONEncoder(conf)
	} else {
		encType = zapcore.NewConsoleEncoder(conf)
	}

	core := zapcore.NewCore(encType, zapcore.AddSync(file), zapcore.WarnLevel)
	logger := zap.New(core, zap.AddStacktrace(zapcore.PanicLevel))

	return &EventLogger{logger: logger}
}
