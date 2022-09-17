package util

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Server struct {
	Mux *http.ServeMux
	s   *http.Server
}

type Context struct {
	Req *http.Request
	Res http.ResponseWriter
}

func NewServer() *Server {
	return &Server{
		Mux: http.NewServeMux(),
	}
}

func (s *Server) Start(addr string) error {
	s.s = &http.Server{
		Addr:    addr,
		Handler: s.Mux,
	}
	err := s.s.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
	return err
}

func (s *Server) StartTLS(addr string, certFile, keyFile string) error {
	s.s = &http.Server{
		Addr:    addr,
		Handler: s.Mux,
	}
	err := s.s.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (s *Server) Stop() {
	if s.s != nil {
		s.s.Shutdown(context.Background())
	}
}

func (s *Server) Static(prefix string, fs http.FileSystem) *Server {
	prefix = strings.TrimSuffix(prefix, "/")
	s.Mux.Handle(prefix+"/", http.StripPrefix(prefix, http.FileServer(fs)))
	return s
}

func (s *Server) StaticDir(prefix, dir string) *Server {
	prefix = strings.TrimSuffix(prefix, "/")
	return s.Static(prefix+"/", http.Dir(dir))
}

func (s *Server) HandleFunc(pattern string, handler func(ctx *Context)) *Server {
	s.Mux.HandleFunc(pattern, func(res http.ResponseWriter, req *http.Request) {
		handler(&Context{
			Req: req,
			Res: res,
		})
	})
	return s
}

func (s *Server) Get(pattern string, handler func(ctx *Context)) *Server {
	s.HandleFunc(pattern, func(ctx *Context) {
		if ctx.Req.Method != http.MethodGet {
			ctx.SetBadRequest()
			return
		}
		handler(ctx)
	})
	return s
}

func (s *Server) Post(pattern string, handler func(ctx *Context)) *Server {
	s.HandleFunc(pattern, func(ctx *Context) {
		if ctx.Req.Method != http.MethodPost {
			ctx.SetBadRequest()
			return
		}
		handler(ctx)
	})
	return s
}

func (s *Server) Proxy(pattern string, target string) *Server {
	pattern = strings.TrimSuffix(pattern, "/")
	target = strings.TrimSuffix(target, "/")
	targetUrl, _ := url.Parse(target)

	// 第一个处理程序，路径有斜杠，处理以此路径为前缀的所有路由
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	s.Mux.Handle(pattern+"/", http.StripPrefix(pattern, proxy))

	// 第二个处理程序，路径没有斜杠，只处理与此路径相等的路由
	proxy = httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host
		req.URL.Path = targetUrl.Path
	}
	s.Mux.Handle(pattern, proxy)

	// 若不定义两个处理程序，则经测试可能会出现在有斜杠和没斜杠之间反复重定向的问题，
	// 例如：访问 /home 会重定向到 /home/，再访问 /home/ 会重定向到 /home，如此反复。
	// 我暂未找到这种情况出现的原因，只能出此下策。
	return s
}

func (c *Context) RedirectTo(target string) *Context {
	http.Redirect(c.Res, c.Req, target, http.StatusFound)
	return c
}

func (c *Context) ProxyTo(target string) *Context {
	targetUrl, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.Director = func(req *http.Request) {
		req.URL = targetUrl
	}
	proxy.ServeHTTP(c.Res, c.Req)
	return c
}

func (c *Context) SetStatus(code int) *Context {
	c.Res.WriteHeader(code)
	return c
}

func (c *Context) SetOK() *Context {
	c.SetStatus(http.StatusOK)
	return c
}

func (c *Context) SetBadRequest(msg ...string) *Context {
	if len(msg) > 0 {
		c.SendText(strings.Join(msg, "\n"))
	}
	c.SetStatus(http.StatusBadRequest)
	return c
}

func (c *Context) SetUnauthorized(msg ...string) *Context {
	if len(msg) > 0 {
		c.SendText(strings.Join(msg, "\n"))
	}
	c.SetStatus(http.StatusUnauthorized)
	return c
}

func (c *Context) SetForbidden(msg ...string) *Context {
	if len(msg) > 0 {
		c.SendText(strings.Join(msg, "\n"))
	}
	c.SetStatus(http.StatusForbidden)
	return c
}

func (c *Context) SetNotFound() *Context {
	c.SetStatus(http.StatusNotFound)
	return c
}

func (c *Context) SetInternalServerError(msg ...string) *Context {
	if len(msg) > 0 {
		c.SendText(strings.Join(msg, "\n"))
	}
	c.SetStatus(http.StatusInternalServerError)
	return c
}

func (c *Context) GetQuery(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) GetJSON(data any) error {
	err := json.NewDecoder(c.Req.Body).Decode(data)
	if err != nil {
		c.SetBadRequest("请求体解析失败")
	}
	return err
}

func (c *Context) SendJSON(data any) *Context {
	c.Res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c.Res).Encode(data)
	return c
}

func (c *Context) SendText(text string) *Context {
	c.Res.Header().Set("Content-Type", "text/plain")
	c.Res.Write([]byte(text))
	return c
}

func (c *Context) SendFileForPath(path string) *Context {
	http.ServeFile(c.Res, c.Req, path)
	return c
}

func (c *Context) SendFileForFS(fs http.FileSystem, path string) *Context {
	file, err := fs.Open(path)
	if err != nil {
		c.SetNotFound()
		return c
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		c.SetInternalServerError()
		return c
	}
	http.ServeContent(c.Res, c.Req, path, stat.ModTime(), file)
	return c
}
