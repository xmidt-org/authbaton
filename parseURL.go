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
