package view

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"proxy/config"
	"proxy/lib"
)

func StartViewServer(dist *embed.FS, proxy *lib.Proxy) {
	config := config.ReadConfig()
	mux := http.NewServeMux()

	distFs, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Println(err)
		return
	}

	fs := http.FileServer(http.FS(distFs))
	mux.Handle("/", fs)

	handleGetConfig(mux)
	handleSetConfig(mux, proxy)

	server := http.Server{
		Addr:    ":9000",
		Handler: mux,
	}

	log.Println("view server start at http://localhost:9000")

	if config.HTTPS != nil {
		err := server.ListenAndServeTLS(config.HTTPS.CertFile, config.HTTPS.KeyFile)
		if err != nil {
			log.Println(err)
		}
	} else {
		err := server.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}
}

func handleGetConfig(mux *http.ServeMux) {
	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		log.Println("config")
	})
}

func handleSetConfig(mux *http.ServeMux, proxy *lib.Proxy) {
	mux.HandleFunc("/api/configSet", func(w http.ResponseWriter, r *http.Request) {
		log.Println("configSet")
		proxy.StartProxyServer()
	})
}
