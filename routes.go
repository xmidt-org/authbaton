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
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/xmidt-org/httpaux"
	"github.com/xmidt-org/touchstone/touchhttp"
	"go.uber.org/fx"
)

type PrimaryRouterIn struct {
	fx.In
	Router    *mux.Router `name:"server_primary"`
	AuthChain alice.Chain `name:"primary_auth_chain"`
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

func buildPrimaryRoutes(in PrimaryRouterIn) {
	in.Router.Use(in.AuthChain.Then)
	in.Router.Handle(fmt.Sprintf("/%s/auth", apiBase), httpaux.ConstantHandler{StatusCode: http.StatusOK})
}

func buildMetricsRoutes(in MetricsRoutesIn) {
	if in.Handler != nil {
		in.Router.Handle("/metrics", in.Handler).Methods("GET")
	}
}

func metricMiddleware(bundle touchhttp.ServerBundle) (out MetricMiddlewareOut) {
	out.Primary = alice.New(bundle.ForServer("server_primary").Then)
	out.Health = alice.New(bundle.ForServer("server_health").Then)
	return
}
