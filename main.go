package main

import (
	"embed"
	"proxy/lib"
)

//go:embed dist
var dist embed.FS

func main() {
	proxy := &lib.Proxy{
		Logger: lib.NewLogger(),
	}
	config := &lib.Config{}
	config.Get()
	go proxy.StartProxyServer(config)
	lib.StartViewServer(&dist, proxy, config)
}
