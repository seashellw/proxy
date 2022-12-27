package src

import (
	"log"
	"proxy/lib"
	"strings"
)

type Proxy struct {
	Server *lib.Server
	Config *Config
}

func (proxy *Proxy) Start() {
	if proxy.Server != nil {
		proxy.Stop()
	}

	proxy.Server = lib.NewServer()
	c := proxy.Config

	if c.Service != nil {
		for key, value := range c.Service {
			proxy.Server.Proxy(key, value)
		}
	}

	if c.Static != nil {
		for key, value := range c.Static {
			proxy.Server.StaticDir(key, value)
		}
	}

	if c.Redirect != nil {
		for key, value := range c.Redirect {
			source := strings.TrimSuffix(key, "/") + "/"
			proxy.Server.HandleFunc(source, func(ctx *lib.Context) {
				if ctx.Req.URL.Path == source {
					ctx.RedirectTo(value)
				}
			})
		}
	}

	log.Println("proxy server start")

	if c.HTTPS != nil {
		err := proxy.Server.StartTLS(":443", c.HTTPS.CertFile, c.HTTPS.KeyFile)
		if err != nil {
			log.Println(err)
		}
	} else {
		err := proxy.Server.Start(":80")
		if err != nil {
			log.Println(err)
		}
	}
}

func (proxy *Proxy) Stop() {
	proxy.Server.Stop()
	log.Println("proxy server stop")
}
