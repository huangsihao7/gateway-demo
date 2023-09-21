package main

import (
	"context"
	"fmt"
	"gateway-detail/tcp/tcp_details"
	"log"
	"net"
)

var (
	addr = ":7002"
)

type TcpHandler struct {
}

func (t *TcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("I am tcpHandler\n"))
}

func main() {
	//tcp服务器测试
	log.Println("Starting tcpserver at " + addr)
	tcpServ := tcp_details.TcpServer{
		Addr:    addr,
		Handler: &TcpHandler{},
	}
	fmt.Println("Starting tcp_server at " + addr)
	log.Fatal(tcpServ.ListenAndServe())

}
