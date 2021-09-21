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

func metricMiddleware(bundle touchhttp.ServerBundle) (out MetricMiddlewareOut) {
	out.Primary = alice.New(bundle.ForServer("server_primary").Then)
	out.Health = alice.New(bundle.ForServer("server_health").Then)
	return
}

func provideServers() fx.Option {
	return fx.Options(
		fx.Provide(
			metricMiddleware,
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
