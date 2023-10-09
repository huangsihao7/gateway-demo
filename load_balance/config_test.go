package main

import (
	"fmt"
	"testing"
)

func TestNewLoadBalanceObserver(t *testing.T) {
	moduleConf, err := NewLoadBalanceZkConf("%s", "/real_server",
		[]string{"172.22.110.88:2181"}, map[string]string{"127.0.0.1:2003": "20"})
	if err != nil {
		fmt.Println("err", err)
		return
	}
	// 为什么要把自己变成观察者 然后在加入自己呢
	loadBalanceObserver := NewLoadBalanceObserver(moduleConf)
	moduleConf.Attach(loadBalanceObserver)
	moduleConf.UpdateConf([]string{"122.11.11"})
	select {}

}
