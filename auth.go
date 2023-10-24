// SPDX-FileCopyrightText: 2021 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/xmidt-org/arrange"
	"github.com/xmidt-org/bascule/basculechecks"
	"github.com/xmidt-org/bascule/basculehttp"
	"go.uber.org/fx"
)

func provideAuth(configKey string) fx.Option {
	return fx.Options(
		basculehttp.ProvideMetrics(),
		basculechecks.ProvideMetrics(),
		fx.Provide(
			arrange.UnmarshalKey("onErrorHTTPResponse", onErrorHTTPResponseConfig{AuthType: "Bearer"}),
			arrange.UnmarshalKey("parseURL", parseURLConfig{URLPathPrefix: "/"}),
			onErrorHTTPResponse,
			parseURLFunc,
		),
		basculehttp.ProvideBasicAuth(configKey),
		basculehttp.ProvideBearerTokenFactory(fmt.Sprintf("%s.bearer", configKey), true),
		basculechecks.ProvideRegexCapabilitiesValidator(fmt.Sprintf("%s.capabilities", configKey)),
		basculehttp.ProvideBearerValidator(),
		basculehttp.ProvideServerChain(),
	)

}
