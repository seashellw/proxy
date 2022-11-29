package lib

import (
	"log"
	"proxy/hp"
	"strings"
)

type Proxy struct {
	Server *hp.Server
	Config *Config
}

func (proxy *Proxy) Start() {
	if proxy.Server != nil {
		proxy.Stop()
	}

	proxy.Server = hp.NewServer()
	c := proxy.Config

	if c.Service != nil {
		for _, service := range c.Service {
			proxy.Server.Proxy(service.Path, service.Target)
		}
	}

	if c.Static != nil {
		for _, service := range c.Static {
			proxy.Server.StaticDir(service.Path, service.Dir)
		}
	}

	if c.Redirect != nil {
		for _, service := range c.Redirect {
			service.Path = strings.TrimSuffix(service.Path, "/")
			proxy.Server.HandleFunc(service.Path+"/", func(ctx *hp.Context) {
				if ctx.Req.URL.Path == service.Path+"/" {
					ctx.RedirectTo(service.Target)
				}
			})
		}
	}

	if c.DynamicService != nil {
		service := *c.DynamicService
		proxy.Server.HandleFunc(service.Path, func(ctx *hp.Context) {
			target := ctx.GetQuery(service.Query)
			if target == "" {
				ctx.SetBadRequest()
				return
			}
			ctx.ProxyTo(target)
		})
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
