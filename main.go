/**
 * Copyright 2021 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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
