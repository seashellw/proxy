package lib

import (
	"proxy/util"
	"strings"
)

type Proxy struct {
	Server *util.Server
	Logger *Logger
}

func (proxy *Proxy) StartProxyServer(config *Config) {
	if proxy.Server != nil {
		proxy.StopProxyServer()
	}

	proxy.Server = util.NewServer()

	if config.Service != nil {
		for _, service := range config.Service {
			proxy.Server.Proxy(service.Path, service.Target)
		}
	}

	if config.Static != nil {
		for _, service := range config.Static {
			proxy.Server.StaticDir(service.Path, service.Dir)
		}
	}

	if config.Redirect != nil {
		for _, service := range config.Redirect {
			service.Path = strings.TrimSuffix(service.Path, "/")
			proxy.Server.HandleFunc(service.Path+"/", func(ctx *util.Context) {
				if ctx.Req.URL.Path == service.Path+"/" {
					ctx.RedirectTo(service.Target)
				}
			})
		}
	}

	if config.DynamicService != nil {
		service := *config.DynamicService
		proxy.Server.HandleFunc(service.Path, func(ctx *util.Context) {
			target := ctx.GetQuery(service.Query)
			if target == "" {
				ctx.SetBadRequest("代理目标为空")
				return
			}
			ctx.ProxyTo(target)
		})
	}

	proxy.Logger.Info("proxy server start")

	if config.HTTPS != nil {
		err := proxy.Server.StartTLS(":443", config.HTTPS.CertFile, config.HTTPS.KeyFile)
		if err != nil {
			proxy.Logger.Error(err.Error())
		}
	} else {
		err := proxy.Server.Start(":80")
		if err != nil {
			proxy.Logger.Error(err.Error())
		}
	}
}

func (proxy *Proxy) StopProxyServer() {
	proxy.Server.Stop()
	proxy.Logger.Info("proxy server stop")
}
