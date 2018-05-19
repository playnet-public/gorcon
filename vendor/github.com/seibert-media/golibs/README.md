# //S/M Go Libs
[![Go Report Card](https://goreportcard.com/badge/github.com/seibert-media/golibs)](https://goreportcard.com/report/github.com/seibert-media/golibs)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Build Status](https://travis-ci.org/seibert-media/golibs.svg?branch=master)](https://travis-ci.org/seibert-media/golibs)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/f61779459d564fb59fc1013d27b36b1f)](https://www.codacy.com/app/seibert-media/golibs?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=seibert-media/golibs&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://api.codacy.com/project/badge/Coverage/f61779459d564fb59fc1013d27b36b1f)](https://www.codacy.com/app/seibert-media/golibs?utm_source=github.com&utm_medium=referral&utm_content=seibert-media/golibs&utm_campaign=Badge_Coverage)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/seibert-media/golibs)

The repository contains various shared libs for use in //SEIBERT/MEDIA Golang projects.

## Libs

### Logging
This logging setup is using go.uber.org/zap.
Sentry is being added for production environments.

```go
l := log.New(
    "sentryDSN",
    false,
)
```

Afterwards the logger can be used just like a default zap.Logger.
When the log level is Error or worse, a sentry message is being sent containing all string and int tags.
If you provide a zap.Error tag, the related stacktrace will also be attached.

To directly access Sentry the internal client is public.

The new implementation found here implements the `context.Context` interface and can therefor be used as a drop in replacement (in most cases).

To do so, just replace your usual `context` import with `context "github.com/seibert-media/golibs/log"`.

As we don't have full compatibility with the native context helper functions like `context.WithValue(...)` this implementation is
providing built-in functions to achieve the same.
Namely:
- `(c *Context) WithValue(key, val interface{}) *Context`
- `(c *Context) WithCancel() (*Context, context.CancelFunc)`
- `(c *Context) WithDeadline(d time.Time) (*Context, context.CancelFunc)`
- `(c *Context) WithTimeout(d time.Duration) (*Context, context.CancelFunc)`
all of this can be accessed directly through the object to be modified.

Aside from that, the original implementation has been migrated into this package to allow replacing
the original context import with this package and provide still working code.
Namely:
- `WithValue(c *Context, key, val interface{}) *Context`
- `WithCancel(c *Context) (*Context, context.CancelFunc)`
- `WithDeadline(c *Context, d time.Time) (*Context, context.CancelFunc)`
- `WithTimeout(c *Context, d time.Duration) (*Context, context.CancelFunc)`

There might still be incompatibility issues if dependencies do internal context modification and do not use this import.

## Compatibility

This library requires Go 1.9+ and is currently tested against Go 1.9.x and 1.10.x
For an up-to-date status on this check [.travis.yml](.travis.yml).

## Contributions

Pull Requests and Issue Reports are welcome.
