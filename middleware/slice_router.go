package middleware

import (
	"context"
	"math"
	"net/http"
)

const abortIndex int8 = math.MaxInt8 / 2 //最多 63 个中间件

type HandlerFunc func(*SliceRouterContext)

// SliceRouter router 结构体
type SliceRouter struct {
	groups []*SliceGroup
}

// SliceGroup group 结构体
type SliceGroup struct {
	//这个路由组是属于哪个路由的
	*SliceRouter
	path     string
	handlers []HandlerFunc
}

// SliceRouterContext router 上下文
type SliceRouterContext struct {
	Rw  http.ResponseWriter
	Req *http.Request
	Ctx context.Context
	*SliceGroup
	index int8
}

// Reset 重置回调
func (c *SliceRouterContext) Reset() {
	c.index = -1
}

// NewSliceRouter 构建router
func NewSliceRouter() *SliceRouter {
	return &SliceRouter{}
}

// Group 创建组
func (g *SliceRouter) Group(path string) *SliceGroup {
	return &SliceGroup{
		path:        path,
		SliceRouter: g,
	}
}

// 构建 中间件的方法
func (g *SliceGroup) Use(middlewares ...HandlerFunc) *SliceGroup {
	g.handlers = append(g.handlers, middlewares...)
	existsFlag := false

	//看下这个路由组有没有在路由下面 没有就加入
	for _, oldGroup := range g.SliceRouter.groups {
		if oldGroup == g {
			existsFlag = true
		}
	}

	if !existsFlag {
		g.SliceRouter.groups = append(g.SliceRouter.groups, g)
	}
	return g
}

// Next 执行中间件
func (c *SliceRouterContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// Abort 跳出中间件
func (c *SliceRouterContext) Abort() {
	c.index = abortIndex
}

// IsAborted 是否跳出中间件
func (c *SliceRouterContext) IsAborted() bool {
	return c.index >= abortIndex
}
