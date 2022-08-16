package main

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Target   string
	Path     string
	Port     string
	CertFile string
	KeyFile  string
}

func readConfig() *Config {
	file, _ := os.Open("./config.json")
	defer file.Close()
	var config Config
	decoder := json.NewDecoder(file)
	decoder.Decode(&config)
	return &config
}

func main() {
	config := readConfig()
	if !strings.HasSuffix(config.Target, "/") {
		config.Target = config.Target + "/"
	}
	if !strings.HasSuffix(config.Path, "/") {
		config.Path = config.Path + "/"
	}
	target, _ := url.Parse(config.Target)
	proxy := httputil.NewSingleHostReverseProxy(target)
	http.Handle(config.Path, http.StripPrefix(config.Path, proxy))
	http.ListenAndServeTLS(":"+config.Port, config.CertFile, config.KeyFile, nil)
}
