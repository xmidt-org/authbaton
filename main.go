// SPDX-FileCopyrightText: 2021 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/xmidt-org/arrange"
	"github.com/xmidt-org/touchstone"

	"github.com/xmidt-org/touchstone/touchhttp"

	"github.com/spf13/pflag"

	"go.uber.org/fx"
)

const (
	applicationName = "authbaton"
	defaultKeyID    = "current"
)

var (
	GitCommit = "undefined"
	Version   = "undefined"
	BuildTime = "undefined"
)

func main() {
	v, logger, err := setup(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	app := fx.New(
		arrange.LoggerFunc(logger.Sugar().Infof),
		arrange.ForViper(v),
		fx.Supply(logger),
		fx.Supply(v),
		provideAuth("authx.inbound"),
		touchstone.Provide(),
		touchhttp.Provide(),
		fx.Provide(
			consts,
			arrange.UnmarshalKey("prometheus", touchstone.Config{}),
			arrange.UnmarshalKey("prometheus.handler", touchhttp.Config{}),
		),
		provideServers(),
	)

	switch err := app.Err(); {
	case errors.Is(err, pflag.ErrHelp):
		return
	case err == nil:
		app.Run()
	default:
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

// Provide the constants in the main package for other uber fx components to use.
type ConstOut struct {
	fx.Out
	DefaultKeyID string `name:"default_key_id"`
}

func consts() ConstOut {
	return ConstOut{
		DefaultKeyID: defaultKeyID,
	}
}
