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
			res, err := http.Get(cdn + path)
			if err != nil {
				log.Println(err)
				continue
			}
			if res.StatusCode != http.StatusOK {
				log.Println(res.Status, cdn+path)
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
