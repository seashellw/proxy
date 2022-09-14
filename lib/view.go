package lib

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"strconv"
)

type ListResponse struct {
	List [][]string
}

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
		password := r.URL.Query().Get("password")
		if password != config.Password {
			proxy.Logger.Info([]string{"/api/config", "password error, password: " + password})
			w.WriteHeader(http.StatusForbidden)
			return
		}
		resConfig := *config
		resConfig.Password = ""

		bytes, err := json.Marshal(resConfig)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	})

	mux.HandleFunc("/api/configSet", func(w http.ResponseWriter, r *http.Request) {
		password := r.URL.Query().Get("password")
		req := &Config{}
		json.NewDecoder(r.Body).Decode(req)
		req.Password = password
		// 读取最新的密码
		config.Read()

		if req.Password != config.Password {
			proxy.Logger.Info([]string{"/api/configSet", "password error, password: " + req.Password})
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		bytes, err := json.Marshal(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		config.Write(bytes)
		proxy.Logger.Info([]string{"config set", string(bytes)})
		go proxy.StartProxyServer(config)
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/api/log", func(w http.ResponseWriter, r *http.Request) {
		password := r.URL.Query().Get("password")
		if password != config.Password {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		start, err := strconv.Atoi(r.URL.Query().Get("start"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		end, err := strconv.Atoi(r.URL.Query().Get("end"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res := ListResponse{
			List: proxy.Logger.Read(start, end),
		}

		jsonStr, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonStr)
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
			proxy.Logger.Error([]string{err.Error()})
		}
	} else {
		err := server.ListenAndServe()
		if err != nil {
			log.Println(err)
			proxy.Logger.Error([]string{err.Error()})
		}
	}
}
