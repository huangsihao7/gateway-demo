package random

import (
	"errors"
	"math/rand"
)

type RandomBalance struct {
	//当前索引
	curIndex int
	//存储负载均衡的地址
	rss []string
	//观察主体
	//conf LoadBalanceConf
}

func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

// Next 随机到下一个index
func (r *RandomBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]
}

// Get 拿到下一个地址
func (r *RandomBalance) Get(key string) (string, error) {
	return r.Next(), nil
}
