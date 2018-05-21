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

// Context is meant to replace context.Context while adding our own logging and tracing functionality
type Context interface {
	context.Context

	WithFields(fields ...zapcore.Field) Context
	WithValue(key, val interface{}) Context
	WithCancel() (Context, context.CancelFunc)
	WithDeadline(d time.Time) (Context, context.CancelFunc)
	WithTimeout(d time.Duration) (Context, context.CancelFunc)

	Ctx() context.Context
	SetCtx(to context.Context)
	Log() *zap.Logger
	SetLogger(to *zap.Logger)
	Sentry() *raven.Client

	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Sync() error
}

// Logger implements Context
type Logger struct {
	context.Context
	*zap.Logger
	SentryClient *raven.Client

	dsn   string
	debug bool
}

// New Context with included logger and sentry instances
// This initializes a new empty context using context.Background()
// To use a custom context call NewWithCustom(...)
func New(dsn string, debug bool) Context {
	sentry, err := raven.New(dsn)
	if err != nil {
		panic(err)
	}

	logger := buildLogger(sentry, debug)

	return &Logger{
		Context:      context.Background(),
		Logger:       logger,
		SentryClient: sentry,

		dsn:   dsn,
		debug: debug,
	}
}

// NewWithCustom Context does what New(...) does but uses the provided context
// instead of initializing a new one
func NewWithCustom(ctx context.Context, dsn string, debug bool) Context {
	sentry, err := raven.New(dsn)
	if err != nil {
		panic(err)
	}

	logger := buildLogger(sentry, debug)

	return &Logger{
		Context:      ctx,
		Logger:       logger,
		SentryClient: sentry,

		dsn:   dsn,
		debug: debug,
	}
}

// Background does what context.Background would do
// but initializes empty logger and sentry clients
func Background() Context {
	sentry, _ := raven.New("")
	logger := zap.NewNop()

	log := &Logger{
		Context:      context.Background(),
		Logger:       logger,
		SentryClient: sentry,
	}

	return log
}

// NewNop returns Context with empty logging and tracing
func NewNop() Context {
	sentry, _ := raven.New("")
	logger := zap.NewNop()

	log := &Logger{
		Context:      context.Background(),
		Logger:       logger,
		SentryClient: sentry,
	}

	return log
}

// Ctx returns the actual context inside Logger
func (c *Logger) Ctx() context.Context {
	return c.Context
}

// SetCtx overwrites the actual context inside Logger
func (c *Logger) SetCtx(to context.Context) {
	c.Context = to
}

// Log returns the actual Logger
func (c *Logger) Log() *zap.Logger {
	return c.Logger
}

// SetLogger overwrites the actual Logger
func (c *Logger) SetLogger(to *zap.Logger) {
	c.Logger = to
}

// Sentry returns the actual sentry client inside Logger
func (c *Logger) Sentry() *raven.Client {
	return c.SentryClient
}

// WithFields wrapper around zap.With
func (c *Logger) WithFields(fields ...zapcore.Field) Context {
	l := NewWithCustom(c.Ctx(), c.dsn, c.debug)
	l.SetLogger(c.Log().With(fields...))
	return l
}

// WithValue is meant to replace context.WithValue as we can not provide compatibility with it
func (c *Logger) WithValue(key, val interface{}) Context {
	c.SetCtx(context.WithValue(c.Ctx(), key, val))
	return c
}

// WithCancel is meant to replace context.WithCancel as we can not provide compatibility with it
func (c *Logger) WithCancel() (Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(c.Ctx())
	c.SetCtx(ctx)
	return c, cancel
}

// WithDeadline is meant to replace context.WithDeadline as we can not provide compatibility with it
func (c *Logger) WithDeadline(d time.Time) (Context, context.CancelFunc) {
	ctx, cancel := context.WithDeadline(c.Ctx(), d)
	c.SetCtx(ctx)
	return c, cancel
}

// WithTimeout is meant to replace context.WithTimeout as we can not provide compatibility with it
func (c *Logger) WithTimeout(d time.Duration) (Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Ctx(), d)
	c.SetCtx(ctx)
	return c, cancel
}

// WithValue is meant to replace context.WithValue
func WithValue(c Context, key, val interface{}) Context {
	c.SetCtx(context.WithValue(c.Ctx(), key, val))
	return c
}

// WithCancel is meant to replace context.WithCancel
func WithCancel(c Context) (Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(c.Ctx())
	c.SetCtx(ctx)
	return c, cancel
}

// WithDeadline is meant to replace context.WithDeadline
func WithDeadline(c Context, d time.Time) (Context, context.CancelFunc) {
	ctx, cancel := context.WithDeadline(c.Ctx(), d)
	c.SetCtx(ctx)
	return c, cancel
}

// WithTimeout is meant to replace context.WithTimeout
func WithTimeout(c Context, d time.Duration) (Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Ctx(), d)
	c.SetCtx(ctx)
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
