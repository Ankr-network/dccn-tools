package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Ankr-network/dccn-tools/metadata"
	"github.com/Ankr-network/dccn-tools/snowflake"
	"github.com/Ankr-network/dccn-tools/zap"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	Flush() error
	Generate() snowflake.ID
}

const (
	TraceID      = "traceID"
	ParentSpanID = "parentSpanID"
	SpanID       = "spanID"
	HostName     = "HostName"
	startTime    = 1448587470 // the world begin time
)

type handler struct {
	logger *zap.Logger
	node   *snowflake.Node
}

func (h *handler) Generate() snowflake.ID {
	return h.node.Generate()
}

func (h *handler) Flush() error {
	return h.logger.Sync()
}

func (h *handler) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Debug(msg, fields...)
	} else if md, ok := metadata.FromOutgoingContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Debug(msg, fields...)
	}
}

func (h *handler) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Info(msg, fields...)
	} else if md, ok := metadata.FromOutgoingContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Info(msg, fields...)
	}
}

func (h *handler) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Warn(msg, fields...)
	} else if md, ok := metadata.FromOutgoingContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Warn(msg, fields...)
	}
}

func (h *handler) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Error(msg, fields...)
	} else if md, ok := metadata.FromOutgoingContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Error(msg, fields...)
	}
}

func (h *handler) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fields = appendFields(md, fields...)
		h.logger.Fatal(msg, fields...)
	} else if md, ok := metadata.FromOutgoingContext(ctx); ok {
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

func appendFields(md metadata.MD, fields ...zap.Field) []zap.Field {
	if vs := md.Get(TraceID); len(vs) > 0 {
		fields = append(fields, zap.String(TraceID, vs[len(vs)-1]))
	} else {
		fmt.Println("no trace id")
	}
	if vs := md.Get(ParentSpanID); len(vs) > 0 {
		fields = append(fields, zap.String(ParentSpanID, vs[len(vs)-1]))
	}
	if vs := md.Get(SpanID); len(vs) > 0 {
		fields = append(fields, zap.String(SpanID, vs[len(vs)-1]))
	}
	hostName, _ := os.Hostname()
	fields = append(fields, zap.String(HostName, hostName))
	return fields
}
