package hp

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Context struct {
	Req *http.Request
	Res http.ResponseWriter
}

// RedirectTo 重定向
func (c *Context) RedirectTo(target string) *Context {
	http.Redirect(c.Res, c.Req, target, http.StatusFound)
	return c
}

// ProxyTo 代理到特定url
func (c *Context) ProxyTo(target string) *Context {
	targetUrl, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy.Director = func(req *http.Request) {
		req.URL = targetUrl
	}
	proxy.ServeHTTP(c.Res, c.Req)
	return c
}

// SetStatus 设置响应状态码
func (c *Context) SetStatus(code int) *Context {
	c.Res.WriteHeader(code)
	return c
}

// SetOK 设置状态码：StatusOK
func (c *Context) SetOK() *Context {
	c.SetStatus(http.StatusOK)
	return c
}

// SetBadRequest 设置状态码：StatusBadRequest，可选自定义响应体
func (c *Context) SetBadRequest(msg ...string) *Context {
	c.SetStatus(http.StatusBadRequest)
	if len(msg) > 0 {
		c.SendText(strings.Join(msg, "\n"))
	}
	return c
}

// SetUnauthorized 设置状态码：StatusNotFound，可选自定义响应体
func (c *Context) SetUnauthorized(msg ...string) *Context {
	c.SetStatus(http.StatusUnauthorized)
	if len(msg) > 0 {
		c.SendText(strings.Join(msg, "\n"))
	}
	return c
}

// SetForbidden 设置状态码：StatusNotFound，可选自定义响应体
func (c *Context) SetForbidden(msg ...string) *Context {
	c.SetStatus(http.StatusForbidden)
	if len(msg) > 0 {
		c.SendText(strings.Join(msg, "\n"))
	}
	return c
}

func (c *Context) SetNotFound() *Context {
	c.SetStatus(http.StatusNotFound)
	return c
}

// SetInternalServerError 设置状态码：StatusInternalServerError，可选自定义响应体
func (c *Context) SetInternalServerError(msg ...string) *Context {
	c.SetStatus(http.StatusInternalServerError)
	if len(msg) > 0 {
		c.SendText(strings.Join(msg, "\n"))
	}
	return c
}

// GetQuery 获取请求query
func (c *Context) GetQuery(key string) string {
	return c.Req.URL.Query().Get(key)
}

// GetJSON 以JSON格式解码请求体
func (c *Context) GetJSON(data any) error {
	err := json.NewDecoder(c.Req.Body).Decode(data)
	if err != nil {
		c.SetBadRequest("请求体解析失败")
	}
	return err
}

// SendJSON 以JSON格式发送响应
func (c *Context) SendJSON(data any) *Context {
	c.Res.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(c.Res).Encode(data)
	return c
}

// SendText 发送文本响应
func (c *Context) SendText(text string) *Context {
	c.Res.Header().Set("Content-Type", "text/plain")
	_, _ = c.Res.Write([]byte(text))
	return c
}

// SendFileFromPath 发送文件，提供文件路径
func (c *Context) SendFileFromPath(path string) *Context {
	http.ServeFile(c.Res, c.Req, path)
	return c
}

// SendFileForFS 发送文件，提供文件系统和路径
func (c *Context) SendFileForFS(fs http.FileSystem, path string) *Context {
	file, err := fs.Open(path)
	if err != nil {
		c.SetNotFound()
		return c
	}
	stat, err := file.Stat()
	if err != nil {
		c.SetInternalServerError()
		return c
	}
	http.ServeContent(c.Res, c.Req, path, stat.ModTime(), file)
	return c
}
