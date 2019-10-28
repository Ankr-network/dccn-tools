package logger

import (
	"context"
	"testing"

	"github.com/Ankr-network/dccn-tools/metadata"
	"go.uber.org/zap"
)

func TestHandler_Info(t *testing.T) {
	l := NewLogger()
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Metadata{
		TraceID:      l.Generate().String(),
		ParentSpanID: l.Generate().String(),
		SpanID:       l.Generate().String(),
	})
	l.Info(ctx, "body", zap.String("hello", "world"))
	ctx = metadata.NewContext(ctx, metadata.Metadata{
		TraceID:      l.Generate().String(),
		ParentSpanID: l.Generate().String(),
		SpanID:       l.Generate().String(),
	})
	l.Info(ctx, "body", zap.String("together", "world"))
}
