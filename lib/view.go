package lib

import (
	"embed"
	"encoding/json"
	"io/fs"
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
		resConfig := *config
		resConfig.Password = ""

		configJson, err := json.Marshal(resConfig)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(configJson)
	})

	mux.HandleFunc("/api/configSet", func(w http.ResponseWriter, r *http.Request) {
		reqConfig := &Config{}
		json.NewDecoder(r.Body).Decode(reqConfig)
		config.ReadConfig()

		if reqConfig.Password != config.Password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		configJson, err := json.Marshal(reqConfig)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		config.WriteConfig(configJson)
		go proxy.StartProxyServer(config)
		w.WriteHeader(http.StatusOK)
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
