package main

import (
	"embed"
	"proxy/lib"
	"proxy/view"
)

//go:embed dist
var dist embed.FS

func main() {
	proxy := lib.Proxy{}
	go proxy.StartProxyServer()
	view.StartViewServer(&dist, &proxy)
}
