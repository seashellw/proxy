package src

import (
	"io"
	"log"
	"net/http"
	"proxy/lib"
	"sync"
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
		wg := sync.WaitGroup{}
		isReturn := false
		for _, cdn := range cdn.Config.CDN {
			wg.Add(1)
			go func(cdn string) {
				res, err := http.Get(cdn + path)
				if res != nil {
					log.Println(res.Status, cdn+path)
				} else {
					log.Println(err, cdn+path)
				}
				if isReturn || err != nil || res.StatusCode != http.StatusOK {
					wg.Done()
					return
				}
				for key, head := range res.Header {
					for _, val := range head {
						ctx.Res.Header().Add(key, val)
					}
				}
				_, _ = io.Copy(ctx.Res, res.Body)
				isReturn = true
				wg.Done()
			}(cdn)
		}
		wg.Wait()
		if !isReturn {
			ctx.SetNotFound()
		}
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
