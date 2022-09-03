package lib

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	Server *http.Server
}

func (proxy *Proxy) StartProxyServer(config *Config) {
	if proxy.Server != nil {
		proxy.StopProxyServer()
	}
	mux := http.NewServeMux()

	if config.Service != nil {
		for _, service := range config.Service {
			target, _ := url.Parse(service.Target)
			proxy := httputil.NewSingleHostReverseProxy(target)
			mux.Handle(service.Path+"/", http.StripPrefix(service.Path, proxy))
		}
	}

	if config.Redirect != nil {
		for _, service := range config.Redirect {
			mux.HandleFunc(service.Path+"/", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, service.Target, http.StatusMovedPermanently)
			})
		}
	}

	if config.FileService != nil {
		for _, service := range config.FileService {
			mux.Handle(service.Path, http.StripPrefix(service.Path, http.FileServer(http.Dir(service.Dir))))
		}
	}

	if config.DynamicService != nil {
		dynamicService := *config.DynamicService
		mux.HandleFunc(dynamicService.Path, func(w http.ResponseWriter, r *http.Request) {
			target, _ := url.Parse(r.URL.Query().Get(dynamicService.Query))
			targetHost := target.Host
			targetScheme := target.Scheme
			reqHost := r.URL.Host
			reqScheme := r.URL.Scheme

			r.URL = target
			r.URL.Host = reqHost
			r.URL.Scheme = reqScheme

			target, _ = url.Parse(targetScheme + "://" + targetHost)
			proxy := httputil.NewSingleHostReverseProxy(target)
			proxy.ServeHTTP(w, r)
		})
	}

	proxy.Server = &http.Server{
		Handler: mux,
	}

	log.Println("proxy server start")

	if config.HTTPS != nil {
		proxy.Server.Addr = ":443"
		err := proxy.Server.ListenAndServeTLS(config.HTTPS.CertFile, config.HTTPS.KeyFile)
		if err != nil {
			log.Println(err)
		}
	} else {
		proxy.Server.Addr = ":80"
		err := proxy.Server.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}
}

func (proxy *Proxy) StopProxyServer() {
	if proxy.Server != nil {
		proxy.Server.Shutdown(context.Background())
		log.Println("proxy server stop")
	}
}
