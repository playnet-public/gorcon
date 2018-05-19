// Copyright 2018 The GoRcon Authors and 'play-net.org' owners. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/kolide/kit/version"
	context "github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	appName = "GoRcon"
	appKey  = "gorcon"
)

var (
	versionInfo = flag.Bool("version", true, "show version info")
	dbg         = flag.Bool("debug", false, "enable debug mode")
	sentryDsn   = flag.String("sentryDsn", "", "sentry dsn key")
)

func main() {
	flag.Parse()

	if *versionInfo {
		v := version.Version()
		fmt.Printf("-- PlayNet %s --\n", appName)
		fmt.Printf(" - version: %s\n", v.Version)
		fmt.Printf("   branch: \t%s\n", v.Branch)
		fmt.Printf("   revision: \t%s\n", v.Revision)
		fmt.Printf("   build date: \t%s\n", v.BuildDate)
		fmt.Printf("   build user: \t%s\n", v.BuildUser)
		fmt.Printf("   go version: \t%s\n", v.GoVersion)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())

	var zapFields []zapcore.Field
	if !*dbg {
		zapFields = []zapcore.Field{
			zap.String("app", appKey),
			zap.String("version", version.Version().Version),
		}
	}

	ctx := context.New(*sentryDsn, *dbg).WithFields(zapFields...)
	defer ctx.Sync()
	ctx.Info("preparing")

	ctx.Info("finished")
}
