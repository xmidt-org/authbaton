// SPDX-FileCopyrightText: 2021 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"net/http"

	"github.com/xmidt-org/bascule/basculehttp"
)

type onErrorHTTPResponseConfig struct {
	AuthType string
}

func onErrorHTTPResponse(config onErrorHTTPResponseConfig) (basculehttp.OnErrorHTTPResponse, error) {
	if config.AuthType != "Bearer" && config.AuthType != "Basic" {
		return nil, fmt.Errorf("invalid auth type '%s': expected Bearer or Basic", config.AuthType)
	}
	return func(w http.ResponseWriter, reason basculehttp.ErrorResponseReason) {
		switch reason {
		case basculehttp.ChecksNotFound, basculehttp.ChecksFailed, basculehttp.ParseFailed:
			w.WriteHeader(http.StatusForbidden)
		default:
			w.Header().Set(basculehttp.AuthTypeHeaderKey, config.AuthType)
			w.WriteHeader(http.StatusUnauthorized)
		}
	}, nil
}
