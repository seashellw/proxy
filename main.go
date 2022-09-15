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
	config.Read()
	go proxy.StartProxyServer(config)
	lib.StartViewServer(&dist, proxy, config)
}
