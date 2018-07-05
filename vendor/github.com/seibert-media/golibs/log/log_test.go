package log_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
)

func Test_NewDebug(t *testing.T) {
	ctx := log.New("", true)
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Logger == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	ctx.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
	ctx = ctx.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Logger == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	ctx.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}

func Test_From(t *testing.T) {
	l := log.New("", true)
	ctx := context.Background()

	ctx = log.WithLogger(ctx, l)
	if log.From(ctx).IsNop() {
		t.Fatal("logger should not be nop")
	}
	log.From(ctx).Debug("test", zap.String("test", "test"))

	ctx = context.Background()
	if !log.From(ctx).IsNop() {
		t.Fatal("logger should be nop")
	}
}

func Test_NewNoDebug(t *testing.T) {
	ctx := log.New("", false)
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Logger == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	ctx.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
	ctx = ctx.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Logger == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	ctx.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}

func Test_NewInvalidSentryURL(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("New() should have panicked")
			}
		}()
		log.New("^", true)
	}()
}

func Test_NewNop(t *testing.T) {
	ctx := log.NewNop()
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Logger == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	ctx = ctx.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Logger == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	ctx.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}
