package middleware

import (
	"context"
	"net/http"
	"strings"
)

// 新建上下文
func newSliceRouterContext(rw http.ResponseWriter, req *http.Request, r *SliceRouter) *SliceRouterContext {

	newSliceGroup := &SliceGroup{}
	//最长URL前缀匹配
	matchUrlLen := 0
	for _, group := range r.groups {
		if strings.HasPrefix(req.RequestURI, group.path) {
			pathLen := len(group.path)
			if pathLen > matchUrlLen {
				matchUrlLen = pathLen
				*newSliceGroup = *group //浅拷贝数组指针
			}
		}
	}

	c := &SliceRouterContext{
		Rw:         rw,
		Req:        req,
		SliceGroup: newSliceGroup,
		Ctx:        req.Context(),
	}

	c.Reset()
	return c
}
func (c *SliceRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

func (c *SliceRouterContext) Set(key, val interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, val)
}
