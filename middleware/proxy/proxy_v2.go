package proxy

import (
	"elevate-hub/db/meta"
	"elevate-hub/db/models"
	"elevate-hub/svc"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type appProxy struct {
	models.Proxy
	*models.APP
}

type authProxy map[string]*appProxy

func (a authProxy) SortedKeys() []string {
	keys := make([]string, 0, len(a))
	for k := range a {
		keys = append(keys, k)
	}
	// 按长度降序排序，确保最长前缀优先匹配
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})
	return keys
}

type Proxy struct {
	proxys     authProxy
	sortedKeys []string
	mu         sync.RWMutex
	ctx        *svc.ServiceContext
}

func NewProxy(ctx *svc.ServiceContext) *Proxy {
	p := &Proxy{
		proxys: make(authProxy),
		ctx:    ctx,
	}
	p.reloadProxys()
	go p.autoReload(5 * time.Minute)
	return p
}

func (p *Proxy) reloadProxys() {
	apps, _, err := p.ctx.Mysql.APPModel.List(meta.EmptyListOption)
	if err != nil {
		//p.ctx.Logger.Error("Failed to reload proxies:", err)
		return
	}

	newProxys := make(authProxy)
	for _, app := range apps {
		for _, proxy := range app.GetProxy() {
			if _, err := url.Parse(proxy.Addr); err != nil {
				//p.ctx.Logger.Errorf("Invalid proxy address %s: %v", proxy.Addr, err)
				continue
			}
			newProxys[proxy.Prefix] = &appProxy{
				Proxy: proxy,
				APP:   app,
			}
		}
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.proxys = newProxys
	p.sortedKeys = newProxys.SortedKeys()
}

func (p *Proxy) autoReload(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		p.reloadProxys()
	}
}

func (p *Proxy) ServeHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		//p.mu.RLock()
		//defer p.mu.RUnlock()

		if len(p.proxys) == 0 {
			c.Next()
			return
		}

		path := c.Request.URL.Path
		prefix, ok := p.matchPrefix(path)
		if !ok {
			c.Next()
			return
		}

		proxyApp := p.proxys[prefix]
		target, _ := url.Parse(proxyApp.Addr) // 已在上游验证过有效性

		proxy := httputil.NewSingleHostReverseProxy(target)
		originalDirector := proxy.Director

		// 自定义请求处理
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			// 路径重写
			req.URL.Path = strings.TrimPrefix(path, prefix)
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
			// 保留查询参数
			req.URL.RawQuery = c.Request.URL.RawQuery
			// 设置Host头
			req.Host = target.Host
			// 添加X-Forwarded头
			req.Header.Set("X-Forwarded-For", c.ClientIP())
			req.Header.Set("X-Forwarded-Host", c.Request.Host)
		}

		// 错误处理
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			//p.ctx.Logger.Errorf("Proxy error: %v", err)
			w.WriteHeader(http.StatusBadGateway)
		}

		// 执行中间件
		//if proxyApp.APP.NeedsAuth && !p.checkAuth(c) {
		if !p.checkAuth(c) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// 执行代理
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func (p *Proxy) matchPrefix(path string) (string, bool) {
	for _, prefix := range p.sortedKeys {
		if strings.HasPrefix(path, prefix) {
			// 确保精确匹配路径分隔符
			if len(path) == len(prefix) || path[len(prefix)] == '/' {
				return prefix, true
			}
		}
	}
	return "", false
}

func (p *Proxy) checkAuth(c *gin.Context) bool {
	// 实现具体的认证逻辑
	// 例如检查JWT token或API Key
	return true
}
