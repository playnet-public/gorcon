package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// Context implements context.Context while adding our own logging and tracing functionality
type Context struct {
	context.Context
	*zap.Logger
	Sentry *raven.Client

	dsn   string
	debug bool
}

// New Context with included logger and sentry instances
// This initializes a new empty context using context.Background()
// To use a custom context call NewWithCustom(...)
func New(dsn string, debug bool) *Context {
	sentry, err := raven.New(dsn)
	if err != nil {
		panic(err)
	}

	logger := buildLogger(sentry, debug)

	return &Context{
		Context: context.Background(),
		Logger:  logger,
		Sentry:  sentry,

		dsn:   dsn,
		debug: debug,
	}
}

// NewWithCustom Context does what New(...) does but uses the provided context
// instead of initializing a new one
func NewWithCustom(ctx context.Context, dsn string, debug bool) *Context {
	sentry, err := raven.New(dsn)
	if err != nil {
		panic(err)
	}

	logger := buildLogger(sentry, debug)

	return &Context{
		Context: ctx,
		Logger:  logger,
		Sentry:  sentry,

		dsn:   dsn,
		debug: debug,
	}
}

// Background does what context.Background would do
// but initializes empty logger and sentry clients
func Background() *Context {
	sentry, _ := raven.New("")
	logger := zap.NewNop()

	log := &Context{
		Context: context.Background(),
		Logger:  logger,
		Sentry:  sentry,
	}

	return log
}

// NewNop returns Context with empty logging and tracing
func NewNop(ctx context.Context) *Context {
	sentry, _ := raven.New("")
	logger := zap.NewNop()

	log := &Context{
		Context: ctx,
		Logger:  logger,
		Sentry:  sentry,
	}

	return log
}

// WithFields wrapper around zap.With
func (c *Context) WithFields(fields ...zapcore.Field) *Context {
	l := NewWithCustom(c.Context, c.dsn, c.debug)
	l.Logger = l.Logger.With(fields...)
	return l
}

// WithValue is meant to replace context.WithValue as we can not provide compatibility with it
func (c *Context) WithValue(key, val interface{}) *Context {
	c.Context = context.WithValue(c.Context, key, val)
	return c
}

// WithCancel is meant to replace context.WithCancel as we can not provide compatibility with it
func (c *Context) WithCancel() (*Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(c.Context)
	c.Context = ctx
	return c, cancel
}

// WithDeadline is meant to replace context.WithDeadline as we can not provide compatibility with it
func (c *Context) WithDeadline(d time.Time) (*Context, context.CancelFunc) {
	ctx, cancel := context.WithDeadline(c.Context, d)
	c.Context = ctx
	return c, cancel
}

// WithTimeout is meant to replace context.WithTimeout as we can not provide compatibility with it
func (c *Context) WithTimeout(d time.Duration) (*Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Context, d)
	c.Context = ctx
	return c, cancel
}

// WithValue is meant to replace context.WithValue
func WithValue(c *Context, key, val interface{}) *Context {
	c.Context = context.WithValue(c.Context, key, val)
	return c
}

// WithCancel is meant to replace context.WithCancel
func WithCancel(c *Context) (*Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(c.Context)
	c.Context = ctx
	return c, cancel
}

// WithDeadline is meant to replace context.WithDeadline
func WithDeadline(c *Context, d time.Time) (*Context, context.CancelFunc) {
	ctx, cancel := context.WithDeadline(c.Context, d)
	c.Context = ctx
	return c, cancel
}

// WithTimeout is meant to replace context.WithTimeout
func WithTimeout(c *Context, d time.Duration) (*Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Context, d)
	c.Context = ctx
	return c, cancel
}

// NewSentryEncoder with dsn
func NewSentryEncoder(client *raven.Client) zapcore.Encoder {
	return newSentryEncoder(client)
}

func newSentryEncoder(client *raven.Client) *sentryEncoder {
	enc := &sentryEncoder{}
	enc.Sentry = client
	return enc
}

type sentryEncoder struct {
	zapcore.ObjectEncoder
	dsn    string
	Sentry *raven.Client
}

// Clone .
func (s *sentryEncoder) Clone() zapcore.Encoder {
	return newSentryEncoder(s.Sentry)
}

// EncodeEntry .
func (s *sentryEncoder) EncodeEntry(e zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf := buffer.NewPool().Get()
	if e.Level == zapcore.ErrorLevel {
		tags := make(map[string]string)
		var err error
		for _, f := range fields {
			var tag string
			switch f.Type {
			case zapcore.StringType:
				tag = f.String
			case zapcore.Int16Type, zapcore.Int32Type, zapcore.Int64Type:
				tag = fmt.Sprintf("%v", f.Integer)
			case zapcore.ErrorType:
				err = f.Interface.(error)
			}
			tags[f.Key] = tag

		}
		if err == nil {
			s.Sentry.CaptureMessage(e.Message, tags)
			return buf, nil
		}
		s.Sentry.CaptureError(errors.Wrap(err, e.Message), tags)
	}
	return buf, nil
}

func (s *sentryEncoder) AddString(key, val string) {
	tags := s.Sentry.Tags
	if tags == nil {
		tags = make(map[string]string)
	}
	tags[key] = val
	s.Sentry.SetTagsContext(tags)
}

func (s *sentryEncoder) AddInt64(key string, val int64) {
	tags := s.Sentry.Tags
	if tags == nil {
		tags = make(map[string]string)
	}
	tags[key] = fmt.Sprint(val)
	s.Sentry.SetTagsContext(tags)
}

// buildLogger
func buildLogger(sentry *raven.Client, debug bool) *zap.Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel && lvl < zapcore.InfoLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	consoleConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)
	sentryEncoder := NewSentryEncoder(sentry)
	var core zapcore.Core
	if debug {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
			zapcore.NewCore(consoleEncoder, consoleDebugging, debugPriority),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
			zapcore.NewCore(sentryEncoder, consoleErrors, highPriority),
		)
	}

	logger := zap.New(core)
	if debug {
		logger = logger.WithOptions(
			zap.AddCaller(),
			zap.AddStacktrace(zap.ErrorLevel),
		)
	} else {
		logger = logger.WithOptions(
			zap.AddStacktrace(zap.FatalLevel),
		)
	}
	return logger
}
