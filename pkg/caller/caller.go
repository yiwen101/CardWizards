package caller

import (
	"fmt"
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/kitex-contrib/registry-nacos/resolver"
	"github.com/yiwen101/CardWizards/pkg/store"
)

var mu sync.RWMutex
var serviceToClientMap map[string]genericclient.Client
var options []func(*store.ServiceMeta) (client.Option, error)

func init() {
	mu = sync.RWMutex{}
	serviceToClientMap = make(map[string]genericclient.Client)
	options = []func(*store.ServiceMeta) (client.Option, error){
		getServiceRegistryOption,
		getServiceLoadBalancerOption,
	}

	store.InfoStore.RegisterServiceMapListener(&serviceChangeHandler{})
	store.InfoStore.RegisterLoadBalanceChoiceListener(&lbChangeHandler{})
}

type serviceChangeHandler struct{}

func (s *serviceChangeHandler) OnStatechanged(data ...interface{}) error {
	isAdd := data[0].(bool)
	meta := data[1].(*store.ServiceMeta)
	if isAdd {
		return mustUpdateClient(meta)
	}
	return deleteClient(meta)
}

type lbChangeHandler struct{}

func (s *lbChangeHandler) OnStatechanged(data ...interface{}) error {
	meta := data[0].(*store.ServiceMeta)
	return mustUpdateClient(meta)
}

type myProvider struct {
	closeOnce sync.Once
	svcs      chan *descriptor.ServiceDescriptor
}

func newMyProvider(svc *descriptor.ServiceDescriptor) (generic.DescriptorProvider, error) {
	p := &myProvider{
		svcs: make(chan *descriptor.ServiceDescriptor, 1), // unblock with buffered channel
	}
	p.svcs <- svc
	return p, nil
}

func (p *myProvider) Provide() <-chan *descriptor.ServiceDescriptor {
	return p.svcs
}

func (p *myProvider) Close() error {
	p.closeOnce.Do(func() {
		close(p.svcs)
	})
	return nil
}

// service name is not file name, but the name of the service in the idl file
func mustUpdateClient(meta *store.ServiceMeta) error {
	mu.Lock()
	defer mu.Unlock()

	decriptor, err := meta.Descriptor.Get()
	if err != nil {
		return err
	}
	p, err := newMyProvider(decriptor)
	if err != nil {
		return fmt.Errorf("error makring myProvider for %s: %s", meta.ServiceName, err.Error())
	}

	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		return err
	}

	opts, err := getOptionsFor(meta)
	if err != nil {
		return err
	}

	client, err := genericclient.NewClient(
		meta.ServiceName,
		g,
		opts...,
	)
	if err != nil {
		return err
	}
	serviceToClientMap[meta.ServiceName] = client

	return nil
}

func deleteClient(meta *store.ServiceMeta) error {
	mu.Lock()
	defer mu.Unlock()
	delete(serviceToClientMap, meta.ServiceName)
	return nil
}

func GetClient(serviceName string) (genericclient.Client, error) {
	mu.RLock()
	defer mu.RUnlock()
	client, ok := serviceToClientMap[serviceName]
	if !ok {
		return nil, fmt.Errorf("client not found for service %s", serviceName)
	}
	return client, nil
}

func getOptionsFor(meta *store.ServiceMeta) ([]client.Option, error) {
	var opts []client.Option
	for _, optionFunc := range options {
		option, err := optionFunc(meta)
		if err != nil {
			return nil, err
		}
		opts = append(opts, option)
	}
	return opts, nil
}

func getServiceRegistryOption(meta *store.ServiceMeta) (client.Option, error) {
	nacosResolver, err := resolver.NewDefaultNacosResolver()
	if err != nil {
		return client.Option{}, err
	}

	return client.WithResolver(nacosResolver), nil
}

func getServiceLoadBalancerOption(meta *store.ServiceMeta) (client.Option, error) {
	lbChoice := meta.LbType
	switch lbChoice {
	case "default":
		return client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), nil
	case "random":
		return client.WithLoadBalancer(loadbalance.NewWeightedRandomBalancer()), nil
	case "roundrobin":
		return client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()), nil
	default:
		return client.Option{}, fmt.Errorf("load balance choice %s not supported", lbChoice)
	}
}
