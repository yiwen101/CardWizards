package clientOption

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
)

func GetServiceLoadBalancerOption(serviceName string) client.Option {
	// todo: enable optioning
	return client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer())
}
