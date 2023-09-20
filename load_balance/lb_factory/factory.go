package lb_factory

import (
	"gateway-detail/load_balance/consist_hash"
	"gateway-detail/load_balance/random"
	"gateway-detail/load_balance/round"
	"gateway-detail/load_balance/weight"
)

type LbType int

const (
	LbRandom LbType = iota
	LbRoundRobin
	LbWeightRoundRobin
	LbConsistentHash
)

// LoadBalanceFactory 负载均衡的工厂模式
func LoadBalanceFactory(lbType LbType) LoadBalance {
	switch lbType {
	case LbRandom:
		return &random.RandomBalance{}
	case LbRoundRobin:
		return &round.RoundRobinBalance{}
	case LbWeightRoundRobin:
		return &weight.WeightRoundRobinBalance{}
	case LbConsistentHash:
		return consist_hash.NewConsistentHashBanlance(10, nil)
	default:
		return &random.RandomBalance{}
	}
}
