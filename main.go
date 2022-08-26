package main

import (
	"embed"
	"proxy/lib"
)

//go:embed dist
var dist embed.FS

func main() {
	proxy := &lib.Proxy{}
	config := &lib.Config{}
	config.ReadConfig()
	go proxy.StartProxyServer(config)
	lib.StartViewServer(&dist, proxy, config)
}
