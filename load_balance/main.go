package main

import (
	"gateway-detail/load_balance/cmd"
	"gateway-detail/load_balance/lb_factory"
	"log"
	"net/http"
)

func main() {
	rb := lb_factory.LoadBalanceFactory(lb_factory.LbConsistentHash)
	if err := rb.Add("http://127.0.0.1:2003/base", "10"); err != nil {
		log.Println(err)
	}
	if err := rb.Add("http://127.0.0.1:2004/base", "50"); err != nil {
		log.Println(err)
	}
	proxy := cmd.NewMultipleHostsReverseProxy(rb)
	log.Println("Starting httpserver at " + cmd.Addr)
	log.Fatal(http.ListenAndServe(cmd.Addr, proxy))
}
