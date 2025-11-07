package logger

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type RequestID string
type LoggerString string

const (
	LoggerKey    LoggerString = "logger"
	RequestIDKey RequestID    = "request_id"
)

type Logger struct {
	l *zap.Logger
}

func New(ctx context.Context, logLevel string) (context.Context, error) {
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if logLevel == "dev" {
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	l, err := loggerCfg.Build()
	if err != nil {
		return ctx, err
	}
	logger := &Logger{l: l}
	ctx = context.WithValue(ctx, LoggerKey, logger)
	return ctx, nil
}
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	if ctx == nil {
		return &Logger{l: zap.NewNop()} // fallback логгер
	}

	loggerVal := ctx.Value(LoggerKey)
	if loggerVal == nil {
		return &Logger{l: zap.NewNop()} // fallback логгер
	}

	logger, ok := loggerVal.(*Logger)
	if !ok || logger == nil {
		return &Logger{l: zap.NewNop()} // fallback логгер
	}
	return logger

}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestIDKey) != nil {
		fields = append(fields, zap.String(string(RequestIDKey), ctx.Value(RequestIDKey).(string)))
	}
	l.l.Info(msg, fields...)
}
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestIDKey) != nil {
		fields = append(fields, zap.String(string(RequestIDKey), ctx.Value(RequestIDKey).(string)))
	}
	l.l.Fatal(msg, fields...)
}
func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestIDKey) != nil {
		fields = append(fields, zap.String(string(RequestIDKey), ctx.Value(RequestIDKey).(string)))
	}
	l.l.Debug(msg, fields...)
}

//	func LoggerMiddleware(ctx context.Context, next func(w http.ResponseWriter, r *http.Request)) {
//		GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("%s %s", r.Method, r.URL))
//		next
//
// }
func LoggerMiddleware(ctx context.Context, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Started %s %s", r.Method, r.URL.Path))
		GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Completed %s in %v", r.URL.Path, time.Since(start)))
		next.ServeHTTP(w, r)
	})
}

// func LoggerInterceptor(ctx context.Context,
//	req any,
//	info *grpc.UnaryServerInfo,
//	handler grpc.UnaryHandler) (any, error) {
//	guid := uuid.New().String()
//
//	ctx = context.WithValue(ctx, RequestIDKey, guid)
//	GetLoggerFromCtx(ctx).Info(ctx,
//		"Request started",
//		zap.Time("request_time", time.Now()),
//	)
//	return handler(ctx, req)
// }
