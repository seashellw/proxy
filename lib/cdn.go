package lib

import (
	"io"
	"log"
	"net/http"
	"proxy/hp"
)

type CDN struct {
	Config  *Config
	Server  *hp.Server
	CdnList []string
}

func (cdn *CDN) Start() {
	cdn.Server = hp.NewServer()
	cdn.Server.HandleFunc("/", func(ctx *hp.Context) {
		path := ctx.Req.URL.Path
		if path == "/" {
			ctx.SendText("cdn server")
			return
		}
		for _, cdn := range cdn.CdnList {
			cdnRes, err := http.Get(cdn + path)
			if err != nil {
				log.Println(err)
				continue
			}
			if cdnRes.StatusCode != http.StatusOK {
				log.Println(cdnRes.Status, cdn+path)
				continue
			}
			ctx.Res.Header().Set("Content-Type", cdnRes.Header.Get("Content-Type"))
			_, _ = io.Copy(ctx.Res, cdnRes.Body)
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
