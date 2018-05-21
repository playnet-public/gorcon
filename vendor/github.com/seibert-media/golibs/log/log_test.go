package log_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
)

func Test_NewDebug(t *testing.T) {
	ctx := log.New("", true)
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Log() == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry() == nil {
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
	if ctx.Log() == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry() == nil {
		t.Fatal("sentry is nil")
	}
	ctx.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}

func Test_NewNoDebug(t *testing.T) {
	ctx := log.New("", false)
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Log() == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry() == nil {
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
	if ctx.Log() == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry() == nil {
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
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("New() should have panicked")
			}
		}()
		log.NewWithCustom(context.Background(), "^", true)
	}()
}

func Test_NewNop(t *testing.T) {
	ctx := log.NewNop()
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Log() == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry() == nil {
		t.Fatal("sentry is nil")
	}
	ctx = ctx.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if ctx == nil {
		t.Fatal("ctx is nil")
	}
	if ctx.Log() == nil {
		t.Fatal("logger is nil")
	}
	if ctx.Sentry() == nil {
		t.Fatal("sentry is nil")
	}
	ctx.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	ctx.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}

type TestCtxKey string

func Test_ContextWorks(t *testing.T) {
	ctx := log.New("", true)
	ctx.WithValue(TestCtxKey("test"), "test")
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	if ctx.Value(TestCtxKey("test")) != "test" {
		t.Fatal("ctx should contain text")
	}
	ctx, cancel := ctx.WithCancel()
	if ctx.Err() != nil {
		t.Fatal("context should not have error")
	}
	cancel()
	if ctx.Err() != context.Canceled {
		t.Fatal("context should be closed")
	}
	ctx = log.New("", true)
	ctx.WithDeadline(time.Now().Add(1 * time.Millisecond))
	select {
	case <-time.After(10 * time.Millisecond):
		t.Fatal("context should be closed after deadline")
	case <-ctx.Done():
		break
	}
	ctx = log.New("", true)
	ctx.WithTimeout(1 * time.Millisecond)
	select {
	case <-time.After(10 * time.Millisecond):
		t.Fatal("context should be closed after deadline")
	case <-ctx.Done():
		break
	}
}

func Test_ContextReplacementWorks(t *testing.T) {
	ctx := log.New("", true)
	ctx = log.WithValue(ctx, TestCtxKey("test"), "test")
	ctx.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	if ctx.Value(TestCtxKey("test")) != "test" {
		t.Fatal("ctx should contain text")
	}
	ctx, cancel := log.WithCancel(ctx)
	if ctx.Err() != nil {
		t.Fatal("context should not have error")
	}
	cancel()
	if ctx.Err() != context.Canceled {
		t.Fatal("context should be closed")
	}
	ctx = log.New("", true)
	ctx, _ = log.WithDeadline(ctx, time.Now().Add(1*time.Millisecond))
	select {
	case <-time.After(10 * time.Millisecond):
		t.Fatal("context should be closed after deadline")
	case <-ctx.Done():
		break
	}
	ctx = log.New("", true)
	ctx, _ = log.WithTimeout(ctx, 1*time.Millisecond)
	select {
	case <-time.After(10 * time.Millisecond):
		t.Fatal("context should be closed after deadline")
	case <-ctx.Done():
		break
	}

	nativeCtx := context.Background()
	newCtx := log.Background()

	nativeCtx = context.WithValue(nativeCtx, TestCtxKey("test"), "test")
	newCtx = log.WithValue(newCtx, TestCtxKey("test"), "test")

	if nativeCtx.Value(TestCtxKey("test")) != newCtx.Value(TestCtxKey("test")) {
		t.Fatal("native and new context mismatch")
	}

	nativeCtx, nativeCancel := context.WithCancel(nativeCtx)
	nativeCancel()
	if nativeCtx.Err() != context.Canceled {
		t.Fatal("context should be closed")
	}
	newCtx, newCancel := log.WithCancel(newCtx)
	newCancel()
	if newCtx.Err() != context.Canceled {
		t.Fatal("context should be closed")
	}
}
