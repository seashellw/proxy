package main

import (
	"proxy/lib"
)

func main() {
	config := &lib.Config{}
	config.Get()
	proxy := &lib.Proxy{
		Config: config,
	}
	cdn := &lib.CDN{
		Config: config,
		CdnList: []string{
			"https://cdn-1259243245.cos.ap-shanghai.myqcloud.com",
			"https://esm.sh",
			"https://cdn.jsdelivr.net",
		},
	}
	go cdn.Start()
	proxy.Start()
}
