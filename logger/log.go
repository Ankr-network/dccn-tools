package logger

import (
	"context"
	"github.com/Ankr-network/dccn-tools/metadata"
	"go.uber.org/zap"
)


type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
}

const (
	TraceID      = "traceID"
	ParentSpanID = "parentSpanID"
	SpanID       = "spanID"
)

type handler struct {
	logger *zap.Logger
}

func (h *handler) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if md, ok := metadata.FromContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Debug(msg, fields...)
	}
}

func (h *handler) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if md, ok := metadata.FromContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Info(msg, fields...)
	}
}

func (h *handler) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if md, ok := metadata.FromContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Warn(msg, fields...)
	}
}

func (h *handler) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if md, ok := metadata.FromContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Error(msg, fields...)
	}
}

func (h *handler) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if md, ok := metadata.FromContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Fatal(msg, fields...)
	}
}

func NewLogger() Logger {
	l, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	return &handler{
		logger: l,
	}
}

func appendFields(md metadata.Metadata, fields ...zap.Field) []zap.Field {
	fields = append(fields,zap.String(TraceID, md[TraceID]))
	fields = append(fields, zap.String(ParentSpanID, md[ParentSpanID]))
	fields = append(fields, zap.String(SpanID, md[SpanID]))
	return fields
}
