package logger

import (
	"context"
	"testing"

	"github.com/Ankr-network/dccn-tools/zap"
	"google.golang.org/grpc/metadata"
)

func TestHandler_Info(t *testing.T) {
	l := NewLogger()
	ctx := context.Background()
	m := make(map[string]string)
	m[TraceID] = l.Generate().String()
	m[ParentSpanID] = l.Generate().String()
	m[SpanID] = l.Generate().String()
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(m))
	l.Info(ctx, "body", zap.String("hello", "world"))
	l.Warn(ctx, "body", zap.String("hello", "warn"))
	l.Error(ctx, "body", zap.String("hello", "Error"))
}
