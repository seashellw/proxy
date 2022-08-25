package lib

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
)

func StartViewServer(dist *embed.FS, proxy *Proxy, config *Config) {
	mux := http.NewServeMux()

	distFs, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Println(err)
		return
	}

	fs := http.FileServer(http.FS(distFs))
	mux.Handle("/", fs)

	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(config.ReadConfigJson())
	})

	mux.HandleFunc("/api/configSet", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		configText, _ := ioutil.ReadAll(r.Body)
		config.WriteConfig(configText)
		proxy.StartProxyServer(config)
	})

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
