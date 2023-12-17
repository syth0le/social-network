package utils

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type Level string

func (l Level) String() string {
	return string(l)
}

func (l *Level) UnmarshalText(text []byte) error {
	strText := string(text)
	switch strText {
	case TraceLevel.String(), DebugLevel.String(), InfoLevel.String(), WarnLevel.String(), ErrorLevel.String(), FatalLevel.String():
		*l = Level(strText)
		return nil
	default:
		return fmt.Errorf("unexpected level: %s", strText)
	}
}

const (
	TraceLevel Level = "trace"
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	FatalLevel Level = "fatal"
)

type Environment string

const (
	Production  Environment = "prod"
	Development Environment = "dev"
)

func LoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Sugar().Infof("http request: %s", r.RequestURI)
			next.ServeHTTP(w, r)
		})
		return fn
	}
}
