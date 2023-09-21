package middleware

import (
	"fmt"
	"golang.org/x/time/rate"
	"log"
)

func TraceLogSliceMW() func(c *SliceRouterContext) {
	return func(c *SliceRouterContext) {
		log.Println("trace_in")
		c.Next()
		log.Println("trace_out")
	}
}

func RateLimiter() func(c *SliceRouterContext) {
	l := rate.NewLimiter(1, 2)
	return func(c *SliceRouterContext) {
		if !l.Allow() {
			c.Rw.Write([]byte(fmt.Sprintf("rate limit:%v,%v", l.Limit(), l.Burst())))
			c.Abort()
			return
		}
		c.Next()
	}
}
