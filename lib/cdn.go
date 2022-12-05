package lib

import (
	"io"
	"log"
	"net/http"
	"proxy/hp"
)

type CDN struct {
	Config *Config
	Server *hp.Server
}

func (cdn *CDN) Start() {
	if cdn.Server != nil {
		cdn.Stop()
	}
	if cdn.Config.CDNList == nil {
		return
	}
	cdn.Server = hp.NewServer()
	cdn.Server.HandleFunc("/", func(ctx *hp.Context) {
		path := ctx.Req.URL.Path
		if path == "/" {
			ctx.SendText("cdn server")
			return
		}
		for _, cdn := range cdn.Config.CDNList {
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
