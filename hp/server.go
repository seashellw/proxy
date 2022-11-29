package hp

import (
	"context"
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

// NewServer 创建一个新的服务器
func NewServer() *Server {
	return &Server{
		Mux: http.NewServeMux(),
	}
}

// Start 启动
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

// StartTLS HTTPS启动
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

// Stop 停机
func (s *Server) Stop() {
	if s.s != nil {
		_ = s.s.Shutdown(context.Background())
	}
}

// Static 静态资源服务，文件系统模式
func (s *Server) Static(prefix string, fs http.FileSystem) *Server {
	prefix = strings.TrimSuffix(prefix, "/")
	handler := http.StripPrefix(prefix, http.FileServer(fs))
	s.Mux.HandleFunc(prefix+"/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, prefix)
		file, err := fs.Open(name)
		if err != nil {
			r.URL.Path = prefix + "/"
		} else {
			_ = file.Close()
		}
		handler.ServeHTTP(w, r)
	})
	if prefix == "" {
		return s
	}
	s.Mux.Handle(prefix, handler)
	return s
}

// StaticDir 静态资源服务，本机目录模式
func (s *Server) StaticDir(prefix, dir string) *Server {
	return s.Static(prefix, http.Dir(dir))
}

// HandleFunc 添加处理程序
func (s *Server) HandleFunc(pattern string, handler func(ctx *Context)) *Server {
	s.Mux.HandleFunc(pattern, func(res http.ResponseWriter, req *http.Request) {
		handler(&Context{
			Req: req,
			Res: res,
		})
	})
	return s
}

// Get GET请求
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

// Post POST请求
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

// Proxy 代理服务
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
	// 我暂未找到这种情况出现的原因
	return s
}
