package main

import "github.com/xmidt-org/bascule/basculehttp"

type parseURLOption struct {
	URLPathPrefix string
}

func parseURLFunc(o parseURLOption) basculehttp.ParseURL {
	return basculehttp.CreateRemovePrefixURLFunc(o.URLPathPrefix, basculehttp.DefaultParseURLFunc)
}
