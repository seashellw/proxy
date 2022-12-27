package main

import (
	"proxy/src"
)

func main() {
	config := &src.Config{}
	config.Get()
	cdn := &src.CDN{
		Config: config,
	}
	go cdn.Start()
	proxy := &src.Proxy{
		Config: config,
	}
	proxy.Start()
}
