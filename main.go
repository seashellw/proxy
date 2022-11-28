package main

import (
	"embed"
	"io/fs"
	"proxy/lib"
)

//go:embed client
var client embed.FS

func main() {
	proxy := &lib.Proxy{
		Logger: lib.NewLogger(),
	}
	config := &lib.Config{}
	config.Get()
	go proxy.StartProxyServer(config)
	files, _ := fs.Sub(client, "client")
	lib.StartViewServer(&files, proxy, config)
}
