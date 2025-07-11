package proxy

//
//import (
//	"elevate-hub/db/meta"
//	"elevate-hub/db/models"
//	"elevate-hub/svc"
//	"github.com/gin-gonic/gin"
//	"net/http/httputil"
//	"net/url"
//)
//
//type appProxy struct {
//	models.Proxy
//	*models.APP
//}
//type authProxy map[string]*appProxy
//
//func (a authProxy) Keys() []string {
//	var keys = make([]string, 0, len(a))
//	for k := range a {
//		keys = append(keys, k)
//	}
//	return keys
//}
//
//type Proxy struct {
//	proxys authProxy
//	match  *Match
//}
//
//func NewProxy(ctx *svc.ServiceContext) *Proxy {
//	apps, _, err := ctx.Mysql.APPModel.List(meta.EmptyListOption)
//	if err != nil {
//		panic(err)
//	}
//
//	var proxys = make(map[string]*appProxy)
//	for _, app := range apps {
//		for _, proxy := range app.GetProxy() {
//			proxys[proxy.Prefix] = &appProxy{
//				proxy,
//				app,
//			}
//		}
//	}
//	return &Proxy{
//		proxys: proxys,
//		match:  NewMatch(),
//	}
//}
//
//func (p *Proxy) ServeHTTP() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if len(p.proxys) == 0 {
//			c.Next()
//			return
//		}
//		prefix, ok := p.match.Match(c.Request.URL.Path, p.proxys.Keys())
//		if !ok {
//			c.Next()
//			return
//		}
//		proxyApp := p.proxys[prefix]
//		parse, _ := url.Parse(proxyApp.Proxy.Addr)
//		proxy := httputil.NewSingleHostReverseProxy(parse)
//		proxy.ServeHTTP(c.Writer, c.Request)
//		c.Abort()
//		return
//	}
//}
