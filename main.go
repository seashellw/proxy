package main

import (
	"proxy/lib"
)

func main() {
	config := &lib.Config{}
	config.Get()
	cdn := &lib.CDN{
		Config: config,
	}
	go cdn.Start()
	proxy := &lib.Proxy{
		Config: config,
	}
	proxy.Start()
}
