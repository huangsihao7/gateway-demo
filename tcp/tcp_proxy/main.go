package main

import (
	"fmt"
	"gateway-detail/load_balance/lb_factory"
	"gateway-detail/tcp/tcp_details"
	"gateway-detail/tcp/tcp_middleware"
	"gateway-detail/tcp/tcp_reverse_proxy"
	"log"
)

var (
	addr = ":7777"
)

func main() {
	rb := lb_factory.LoadBalanceFactory(lb_factory.LbRandom)
	rb.Add("127.0.0.1:7002", "40")

	proxy := tcp_reverse_proxy.NewTcpLoadBalanceReverseProxy(&tcp_middleware.TcpSliceRouterContext{}, rb)
	tcpServ := tcp_details.TcpServer{
		Addr:    addr,
		Handler: proxy,
	}
	fmt.Println("tcp_proxy start at:" + addr)
	log.Fatal(tcpServ.ListenAndServe())

	//redis服务器测试
	//rb := lb_factory.LoadBalanceFactory(lb_factory.LbWeightRoundRobin)
	//rb.Add("172.22.110.88:6388", "40")
	//
	//proxy := tcp_reverse_proxy.NewTcpLoadBalanceReverseProxy(&tcp_middleware.TcpSliceRouterContext{}, rb)
	//tcpServ := tcp_details.TcpServer{Addr: addr, Handler: proxy}
	//fmt.Println("Starting tcp_proxy at " + addr)
	//log.Fatal(tcpServ.ListenAndServe())

}
