package main

import (
	"fmt"
	"gateway-detail/zook/zookeeper"
	"time"
)

func main() {
	zkManager := zookeeper.NewZkManager([]string{"172.22.110.88:2181"})
	zkManager.GetConnect()
	defer zkManager.Close()
	i := 0
	for {
		zkManager.RegisterServerPath("/real_server", fmt.Sprint(i))
		fmt.Println("Register", i)
		time.Sleep(5 * time.Second)
		i++
	}
}
