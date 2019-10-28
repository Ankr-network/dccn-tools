package logger

import (
	"context"
	"time"

	"github.com/Ankr-network/dccn-tools/metadata"
	"github.com/Ankr-network/dccn-tools/snowflake"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	Generate() snowflake.ID
}

const (
	TraceID      = "traceID"
	ParentSpanID = "parentSpanID"
	SpanID       = "spanID"
	startTime    = 1448587470 // the world begin time
)

type handler struct {
	logger *zap.Logger
	node   *snowflake.Node
}

func (h *handler) Generate() snowflake.ID {
	return h.node.Generate()
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
	n, err := snowflake.NewNode((time.Now().Unix() - startTime) % 1024)
	if err != nil {
		panic(err)
	}
	return &handler{
		logger: l,
		node:   n,
	}
}

func appendFields(md metadata.Metadata, fields ...zap.Field) []zap.Field {
	fields = append(fields, zap.String(TraceID, md[TraceID]))
	fields = append(fields, zap.String(ParentSpanID, md[ParentSpanID]))
	fields = append(fields, zap.String(SpanID, md[SpanID]))
	return fields
}
