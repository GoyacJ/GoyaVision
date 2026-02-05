package logger

import (
	"context"
	"testing"
)

func TestWithRequestID(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-123")
	rid := RequestIDFromContext(ctx)
	if rid != "req-123" {
		t.Errorf("expected 'req-123', got '%s'", rid)
	}
}

func TestRequestIDFromContextEmpty(t *testing.T) {
	ctx := context.Background()
	rid := RequestIDFromContext(ctx)
	if rid != "" {
		t.Errorf("expected empty string, got '%s'", rid)
	}
}

func TestFromContext(t *testing.T) {
	ctx := context.Background()
	l := FromContext(ctx)
	if l == nil {
		t.Error("logger should not be nil")
	}

	ctx = WithRequestID(ctx, "req-456")
	l = FromContext(ctx)
	if l == nil {
		t.Error("logger with request_id should not be nil")
	}
}

func TestDefault(t *testing.T) {
	l := Default()
	if l == nil {
		t.Error("default logger should not be nil")
	}
}
