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
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"

	"github.com/xmidt-org/argus/auth"
	"github.com/xmidt-org/argus/store/db/metric"
	"github.com/xmidt-org/arrange"
	"github.com/xmidt-org/arrange/arrangehttp"
	"github.com/xmidt-org/httpaux"
	"github.com/xmidt-org/sallust/sallustkit"
	"github.com/xmidt-org/touchstone"

	"github.com/xmidt-org/touchstone/touchhttp"

	"github.com/spf13/pflag"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	applicationName = "authbaton"
	apiBase         = "api/v1"
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
		metric.ProvideMetrics(),
		auth.Provide("authx.inbound"),
		touchstone.Provide(),
		touchhttp.Provide(),
		fx.Provide(
			consts,
			gokitLogger,
			arrange.UnmarshalKey("prometheus", touchstone.Config{}),
			arrange.UnmarshalKey("prometheus.handler", touchhttp.Config{}),
			arrange.UnmarshalKey("onErrorHTTPResponse", onErrorHTTPResponseConfig{AuthType: "Bearer"}),
			arrange.UnmarshalKey("parseURL", parseURLConfig{URLPathPrefix: "/"}),
			metricMiddleware,
			fx.Annotated{
				Name:   "primary_bascule_on_error_http_response",
				Target: onErrorHTTPResponse,
			},
			fx.Annotated{
				Name:   "primary_bascule_parse_url",
				Target: parseURLFunc,
			},
		),

		arrangehttp.Server{
			Name: "server_primary",
			Key:  "servers.primary",
			Inject: arrange.Inject{
				PrimaryMMIn{},
			},
		}.Provide(),

		arrangehttp.Server{
			Name: "server_health",
			Key:  "servers.health",
			Inject: arrange.Inject{
				HealthMMIn{},
			},
			Invoke: arrange.Invoke{
				func(r *mux.Router) {
					r.Handle("/health", httpaux.ConstantHandler{
						StatusCode: http.StatusOK,
					}).Methods("GET")
				},
			},
		}.Provide(),

		arrangehttp.Server{
			Name: "server_metrics",
			Key:  "servers.metrics",
		}.Provide(),

		fx.Invoke(
			serverValidator{Key: "servers.primary"}.Validate,
			serverValidator{Key: "servers.metrics"}.Validate,
			serverValidator{Key: "servers.health"}.Validate,
			handlePrimaryEndpoint,
			handledMetricEndpoint,
		),
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

func gokitLogger(l *zap.Logger) log.Logger {
	return sallustkit.Logger{
		Zap: l,
	}
}

// Provide the constants in the main package for other uber fx components to use.
type ConstOut struct {
	fx.Out
	APIBase      string `name:"api_base"`
	DefaultKeyID string `name:"default_key_id"`
}

func consts() ConstOut {
	return ConstOut{
		APIBase:      apiBase,
		DefaultKeyID: defaultKeyID,
	}
}
