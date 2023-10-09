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
		conf := fmt.Sprintf("{name:" + fmt.Sprint(i) + "}")
		zkManager.SetPathData("/rs_server_conf", []byte(conf), int32(i))
		time.Sleep(5 * time.Second)
		i++
	}
}
