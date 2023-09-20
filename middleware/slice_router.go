package middleware

import (
	"math"
	"net/http"
)

const abortIndex int8 = math.MaxInt8 / 2 //最多 63 个中间件

type SliceRouter struct {
}

type SliceGroup struct {
	*SliceRouter
	path    string
	handler []http.HandlerFunc
}
