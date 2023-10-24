// SPDX-FileCopyrightText: 2021 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/xmidt-org/bascule/basculehttp"
)

type parseURLConfig struct {
	URLPathPrefix string
}

func parseURLFunc(o parseURLConfig) basculehttp.ParseURL {
	return basculehttp.CreateRemovePrefixURLFunc(o.URLPathPrefix, basculehttp.DefaultParseURLFunc)
}
