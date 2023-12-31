package main

import (
	"fmt"
	"gateway-detail/zook/zookeeper"
)

// LoadBalanceConf 配置主题
type LoadBalanceConf interface {
	// Attach 添加观察者
	Attach()
	GetConf()
	WatchConf()
	UpdateConf()
}

type LoadBalanceZkConf struct {
	observers    []Observer
	path         string
	zkHosts      []string
	confIpWeight map[string]string
	activeList   []string
	format       string
}

func (s *LoadBalanceZkConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}
func (s *LoadBalanceZkConf) NotifyAllObservers() {
	for _, obs := range s.observers {
		obs.Update()
	}
}

func (s *LoadBalanceZkConf) GetConf() []string {
	var confList []string
	for i := range s.activeList {
		weight, ok := s.confIpWeight[s.activeList[i]]
		if !ok {
			weight = "50"
		}
		confList = append(confList, fmt.Sprintf(s.format, s.activeList[i])+","+weight)
	}
	return confList
}

// WatchConf 更新配置时 监听者也更新
func (s *LoadBalanceZkConf) WatchConf() {
	zkManager := zookeeper.NewZkManager(s.zkHosts)
	zkManager.GetConnect()
	chanList, chanErr := zkManager.WatchServerListByPath(s.path)
	go func() {
		defer zkManager.Close()
		for {
			select {
			case changeErr := <-chanErr:
				fmt.Println("changeErr", changeErr)
			case changedList := <-chanList:
				fmt.Println("watch node changed")
				s.UpdateConf(changedList)
			}
		}
	}()
}

func (s *LoadBalanceZkConf) UpdateConf(conf []string) {
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

func NewLoadBalanceZkConf(format, path string, zkHosts []string, conf map[string]string) (*LoadBalanceZkConf, error) {
	zkManager := zookeeper.NewZkManager(zkHosts)
	zkManager.GetConnect()
	defer zkManager.Close()
	zlist, err := zkManager.GetServerListByPath(path)
	if err != nil {
		return nil, err
	}
	mConf := &LoadBalanceZkConf{
		format:       format,
		activeList:   zlist,
		confIpWeight: conf,
		zkHosts:      zkHosts,
		path:         path,
	}
	mConf.WatchConf()
	return mConf, nil
}

//观察者部分

type Observer interface {
	Update()
}

type LoadBalanceObserver struct {
	ModuleConf *LoadBalanceZkConf
}

func (l *LoadBalanceObserver) Update() {
	fmt.Println("Update get conf:", l.ModuleConf.GetConf())
}

func NewLoadBalanceObserver(conf *LoadBalanceZkConf) *LoadBalanceObserver {
	return &LoadBalanceObserver{
		ModuleConf: conf,
	}
}
