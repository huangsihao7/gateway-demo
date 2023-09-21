package tcp_details

import "net"

type tcpKeepAliveListener struct {
	*net.TCPListener
}
