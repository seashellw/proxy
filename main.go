package main

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type HTTPSConfig struct {
	CertFile string
	KeyFile  string
}

type ServiceConfig struct {
	Target string
	Path   string
}

type DynamicServiceConfig struct {
	Path  string
	Query string
}

type Config struct {
	Service        []ServiceConfig
	DynamicService *DynamicServiceConfig
	HTTPS          *HTTPSConfig
}

func readConfig() *Config {
	file, _ := os.Open("./config.json")
	defer file.Close()
	var config Config
	decoder := json.NewDecoder(file)
	decoder.Decode(&config)
	return &config
}

//获取URL的GET参数
func getUrlArg(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}

func main() {
	config := readConfig()

	if config.Service != nil {
		for _, service := range config.Service {
			target, _ := url.Parse(service.Target)
			proxy := httputil.NewSingleHostReverseProxy(target)
			http.Handle(service.Path+"/", http.StripPrefix(service.Path, proxy))
		}
	}

	if config.DynamicService != nil {
		dynamicService := *config.DynamicService
		http.HandleFunc(dynamicService.Path, func(w http.ResponseWriter, r *http.Request) {
			target, _ := url.Parse(getUrlArg(r, dynamicService.Query))
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

	if config.HTTPS != nil {
		http.ListenAndServeTLS(":443", config.HTTPS.CertFile, config.HTTPS.KeyFile, nil)
	} else {
		http.ListenAndServe(":80", nil)
	}
}
