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
