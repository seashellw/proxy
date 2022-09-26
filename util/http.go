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

// 创建一个新的服务器
func NewServer() *Server {
	return &Server{
		Mux: http.NewServeMux(),
	}
}

// 启动
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

// HTTPS启动
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

// 停机
func (s *Server) Stop() {
	if s.s != nil {
		s.s.Shutdown(context.Background())
	}
}

// 静态资源服务，文件系统模式
func (s *Server) Static(prefix string, fs http.FileSystem) *Server {
	prefix = strings.TrimSuffix(prefix, "/")
	handler := http.StripPrefix(prefix, http.FileServer(fs))
	s.Mux.HandleFunc(prefix+"/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, prefix)
		file, err := fs.Open(name)
		if err != nil {
			r.URL.Path = prefix + "/"
		} else {
			file.Close()
		}
		handler.ServeHTTP(w, r)
	})
	if prefix == "" {
		return s
	}
	s.Mux.Handle(prefix, handler)
	return s
}

// 静态资源服务，本机目录模式
func (s *Server) StaticDir(prefix, dir string) *Server {
	return s.Static(prefix, http.Dir(dir))
}

// 添加处理程序
func (s *Server) HandleFunc(pattern string, handler func(ctx *Context)) *Server {
	s.Mux.HandleFunc(pattern, func(res http.ResponseWriter, req *http.Request) {
		handler(&Context{
			Req: req,
			Res: res,
		})
	})
	return s
}

// GET请求
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

// POST请求
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

// 代理服务
func (s *Server) Proxy(pattern string, target string) *Server {
	pattern = strings.TrimSuffix(pattern, "/")
	target = strings.TrimSuffix(target, "/")
	targetUrl, _ := url.Parse(target)

	// 第一个处理程序，路径有斜杠，处理以此路径为前缀的所有路由
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	s.Mux.Handle(pattern+"/", http.StripPrefix(pattern, proxy))

	if pattern == "" {
		return s
	}

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

// -------------------以下为上下文方法-------------------

// 重定向
func (c *Context) RedirectTo(target string) *Context {
	http.Redirect(c.Res, c.Req, target, http.StatusFound)
	return c
}

// 代理到特定url
func (c *Context) ProxyTo(target string) *Context {
	targetUrl, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.Director = func(req *http.Request) {
		req.URL = targetUrl
	}
	proxy.ServeHTTP(c.Res, c.Req)
	return c
}

// 设置响应状态码
func (c *Context) SetStatus(code int) *Context {
	c.Res.WriteHeader(code)
	return c
}

// 设置状态码：StatusOK
func (c *Context) SetOK() *Context {
	c.SetStatus(http.StatusOK)
	return c
}

// 设置状态码：StatusBadRequest，可选自定义响应体
func (c *Context) SetBadRequest(msg ...string) *Context {
	if len(msg) > 0 {
		c.SendText(strings.Join(msg, "\n"))
	}
	c.SetStatus(http.StatusBadRequest)
	return c
}

// 设置状态码：StatusNotFound，可选自定义响应体
func (c *Context) SetUnauthorized(msg ...string) *Context {
	if len(msg) > 0 {
		c.SendText(strings.Join(msg, "\n"))
	}
	c.SetStatus(http.StatusUnauthorized)
	return c
}

// 设置状态码：StatusNotFound，可选自定义响应体
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

// 设置状态码：StatusInternalServerError，可选自定义响应体
func (c *Context) SetInternalServerError(msg ...string) *Context {
	if len(msg) > 0 {
		c.SendText(strings.Join(msg, "\n"))
	}
	c.SetStatus(http.StatusInternalServerError)
	return c
}

// 获取请求query
func (c *Context) GetQuery(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 以JSON格式解码请求体
func (c *Context) GetJSON(data any) error {
	err := json.NewDecoder(c.Req.Body).Decode(data)
	if err != nil {
		c.SetBadRequest("请求体解析失败")
	}
	return err
}

// 以JSON格式发送响应
func (c *Context) SendJSON(data any) *Context {
	c.Res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c.Res).Encode(data)
	return c
}

// 发送文本响应
func (c *Context) SendText(text string) *Context {
	c.Res.Header().Set("Content-Type", "text/plain")
	c.Res.Write([]byte(text))
	return c
}

// 发送文件，提供文件路径
func (c *Context) SendFileFromPath(path string) *Context {
	http.ServeFile(c.Res, c.Req, path)
	return c
}

// 发送文件，提供文件系统和路径
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
