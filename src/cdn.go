package src

import (
	"io"
	"log"
	"net/http"
	"proxy/lib"
)

type CDN struct {
	Config *Config
	Server *lib.Server
}

func (cdn *CDN) Start() {
	if cdn.Server != nil {
		cdn.Stop()
	}
	if cdn.Config.CDN == nil {
		return
	}
	cdn.Server = lib.NewServer()
	cdn.Server.HandleFunc("/", func(ctx *lib.Context) {
		path := ctx.Req.URL.Path
		if path == "/" {
			ctx.SendText("cdn server")
			return
		}
		for _, cdn := range cdn.Config.CDN {
			res, err := http.Get(cdn + path)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(res.Status, cdn+path)
			if res.StatusCode != http.StatusOK {
				continue
			}
			for key, head := range res.Header {
				for _, val := range head {
					ctx.Res.Header().Add(key, val)
				}
			}
			_, _ = io.Copy(ctx.Res, res.Body)
			return
		}
		ctx.SetNotFound()
	})

	log.Println("cdn proxy server start")

	if cdn.Config.HTTPS != nil {
		err := cdn.Server.StartTLS(":9002", cdn.Config.HTTPS.CertFile, cdn.Config.HTTPS.KeyFile)
		if err != nil {
			log.Println(err)
		}
	} else {
		err := cdn.Server.Start(":9002")
		if err != nil {
			log.Println(err)
		}
	}
}

func (cdn *CDN) Stop() {
	cdn.Server.Stop()
	log.Println("cdn proxy server stop")
}
