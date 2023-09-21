package main

import (
	"fmt"
	"gateway-detail/middleware"
	"gateway-detail/middleware/proxy"
	"log"
	"net/http"
	"net/url"
)

var addr = "127.0.0.1:2222"

func main() {

	// 创建一个反向代理的handler
	reverseProxy := func(c *middleware.SliceRouterContext) http.Handler {
		rs1 := "http://127.0.0.1:2003/base"
		url1, err1 := url.Parse(rs1)
		if err1 != nil {
			log.Println(err1)
		}

		rs2 := "http://127.0.0.1:2004/base"
		url2, err2 := url.Parse(rs2)
		if err2 != nil {
			log.Println(err2)
		}

		urls := []*url.URL{url1, url2}
		return proxy.NewMultipleHostsReverseProxy(c, urls)
	}
	//初始化方法数组路由
	sliceRouter := middleware.NewSliceRouter()

	//测试中间件的使用
	sliceRouter.Group("/base").Use(middleware.RateLimiter(), middleware.TraceLogSliceMW(), func(c *middleware.SliceRouterContext) {
		c.Rw.Write([]byte("test func"))
	})

	//反向代理添加中间件
	sliceRouter.Group("/").Use(middleware.TraceLogSliceMW(), func(c *middleware.SliceRouterContext) {
		fmt.Printf("reverseProxy")
		reverseProxy(c).ServeHTTP(c.Rw, c.Req)
	})

	routerHandler := middleware.NewSliceRouterHandler(nil, sliceRouter)
	log.Fatal(http.ListenAndServe(addr, routerHandler))
}
