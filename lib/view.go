package lib

import (
	"io/fs"
	"net/http"
	"proxy/util"
)

type ListResponse struct {
	List [][]string
}

func StartViewServer(files *fs.FS, proxy *Proxy, config *Config) {
	server := util.NewServer()
	server.Static("/", http.FS(*files))

	server.Get("/api/config", func(ctx *util.Context) {
		password := ctx.GetQuery("password")
		if password != config.Password {
			ctx.SetForbidden()
			return
		}
		config.Get()
		res := *config
		res.Password = ""
		ctx.SendJSON(res)
	})

	server.Post("/api/configSet", func(ctx *util.Context) {
		password := ctx.GetQuery("password")
		if password != config.Password {
			ctx.SetForbidden()
			return
		}
		req := &Config{}
		err := ctx.GetJSON(req)
		if err != nil {
			ctx.SetBadRequest(err.Error())
			return
		}
		err = config.Set(req)
		if err != nil {
			ctx.SetBadRequest(err.Error())
			return
		}
		go proxy.StartProxyServer(config)
		ctx.SetOK()
	})

	proxy.Logger.Info("view server start at http://localhost:9000")

	if config.HTTPS != nil {
		err := server.StartTLS(":9000", config.HTTPS.CertFile, config.HTTPS.KeyFile)
		if err != nil {
			proxy.Logger.Error(err.Error())
		}
	} else {
		err := server.Start(":9000")
		if err != nil {
			proxy.Logger.Error(err.Error())
		}
	}
}
