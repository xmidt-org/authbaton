package main

import (
	"fmt"
	"net/http"

	"github.com/xmidt-org/bascule/basculehttp"
)

type OnErrorHTTPResponseOption struct {
	AuthType string
}

func onErrorHTTPResponse(config OnErrorHTTPResponseOption) (basculehttp.OnErrorHTTPResponse, error) {
	if config.AuthType != "Bearer" && config.AuthType != "Basic" {
		return nil, fmt.Errorf("invalid auth type '%s': expected Bearer or Basic", config.AuthType)
	}
	return func(w http.ResponseWriter, reason basculehttp.ErrorResponseReason) {
		switch reason {
		case basculehttp.ChecksNotFound, basculehttp.ChecksFailed:
			w.WriteHeader(http.StatusForbidden)
		default:
			w.Header().Set(basculehttp.AuthTypeHeaderKey, config.AuthType)
			w.WriteHeader(http.StatusUnauthorized)
		}
	}, nil
}
