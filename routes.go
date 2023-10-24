// SPDX-FileCopyrightText: 2021 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/xmidt-org/arrange"
	"github.com/xmidt-org/arrange/arrangehttp"
	"github.com/xmidt-org/httpaux"
	"github.com/xmidt-org/touchstone/touchhttp"
	"go.uber.org/fx"
)

type PrimaryRouterIn struct {
	fx.In
	Router    *mux.Router `name:"server_primary"`
	AuthChain alice.Chain `name:"auth_chain"`
}

type MetricsRoutesIn struct {
	fx.In
	Router  *mux.Router `name:"server_metrics"`
	Handler touchhttp.Handler
}

type PrimaryMMIn struct {
	fx.In
	Primary alice.Chain `name:"middleware_primary_metrics"`
}

type HealthMMIn struct {
	fx.In
	Health alice.Chain `name:"middleware_health_metrics"`
}

type MetricMiddlewareIn struct {
	fx.In
	Primary touchhttp.ServerInstrumenter `name:"instrumenter_primary"`
	Health  touchhttp.ServerInstrumenter `name:"instrumenter_health"`
}
type MetricMiddlewareOut struct {
	fx.Out
	Primary alice.Chain `name:"middleware_primary_metrics"`
	Health  alice.Chain `name:"middleware_health_metrics"`
}

func handlePrimaryEndpoint(in PrimaryRouterIn) {
	in.Router.Use(in.AuthChain.Then)
	in.Router.PathPrefix("/").Handler(httpaux.ConstantHandler{StatusCode: http.StatusOK})
}

func handledMetricEndpoint(in MetricsRoutesIn) {
	in.Router.Handle("/metrics", in.Handler).Methods("GET")
}

func provideMetricMiddleware(in MetricMiddlewareIn) (out MetricMiddlewareOut) {
	out.Primary = alice.New(in.Primary.Then)
	out.Health = alice.New(in.Health.Then)
	return out
}

func provideServers() fx.Option {
	var bundle touchhttp.ServerBundle
	return fx.Options(
		fx.Provide(
			provideMetricMiddleware,
			fx.Annotated{
				Name: "instrumenter_primary",
				Target: bundle.NewInstrumenter(
					touchhttp.ServerLabel, "server_primary",
				),
			},
			fx.Annotated{
				Name: "instrumenter_health",
				Target: bundle.NewInstrumenter(
					touchhttp.ServerLabel, "server_health",
				),
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
}
