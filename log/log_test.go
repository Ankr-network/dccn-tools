package log

import (
	"context"
	"github.com/Ankr-network/dccn-tools/metadata"
	"go.uber.org/zap"
	"testing"
)

func TestHandler_Info(t *testing.T) {
	l := NewLogger()
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Metadata{
		TraceID:      "223344",
		ParentSpanID: ParentSpanID,
		SpanID:       SpanID,
	})
	l.Info(ctx, "body", zap.String("hello", "world"))
}
