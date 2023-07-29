package caller

/*
This package is about offering the right generic client. Its only responsibility is to keep the clients to offer to be
up to date. For example when a user update a service with a new idl file, new nacos cluster name, load balance choice,
add or delete a service, this package should be able to update the client accordingly.

If new functionalities or choice is to be added in the future, for example, more client options, this package should
register a new handelr to the store package accordingly.
*/

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/kitex-contrib/registry-nacos/resolver"
	"github.com/yiwen101/CardWizards/pkg/store"
	"github.com/yiwen101/CardWizards/pkg/utils"
)

var serviceToClientMap *utils.MutexMap[string, *myClient]
var options []func(*store.ServiceMeta) (client.Option, error)

func init() {
	serviceToClientMap = utils.NewMutexMap[string, *myClient]()
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
		return addOrReplaceClient(meta)
	}
	return deleteClient(meta)
}

type lbChangeHandler struct{}

func (s *lbChangeHandler) OnStatechanged(data ...interface{}) error {
	meta := data[0].(*store.ServiceMeta)
	return addOrReplaceClient(meta)
}

/*
Building a generic client require a channel of descriptor. But the provider in the generic package seems to
block infinitely after emiting a single instance of descriptor and hence is not suitable to be store and
reuse in the store. So we decide to store "descriptor keeper" in store and make the my provider wrapper class
to the service descriptor returned from the store.

The alternatice solution of reparsing the file in this package is rejected. That violates single responsibility and
single source of truth principle. Store should be the only source of metaDatas (for exampele, service descriptors),
and other data (for example, client) should be generared using the infomation in the store. So we implement our own
provider as a wrapper class to build the generic client (which require provider, not descriptor as input)
*/
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
	defer p.Close()
	return p.svcs
}

func (p *myProvider) Close() error {
	p.closeOnce.Do(func() {
		close(p.svcs)
	})
	return nil
}

/*
Generic client does not support pointer receiver, and it is too expensive to copy the whole client
each time we make a call; so I make a wraper class here

An alternative is to include the handler logic here to reduce passing parameters on copy via funciton.
But agian, it will mess up with the single responsibility principle. This class should only be
responsible for keeping clients up to date,
*/

type myClient struct {
	client genericclient.Client
}

func (c *myClient) GenericCall(ctx context.Context, method string, request interface{}, callOptions ...callopt.Option) (response interface{}, err error) {
	return c.client.GenericCall(ctx, method, request, callOptions...)
}

// service name is not file name, but the name of the service in the idl file
func addOrReplaceClient(meta *store.ServiceMeta) error {
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
	serviceToClientMap.AddOrReplace(meta.ServiceName, &myClient{client: client})
	return nil
}

func deleteClient(meta *store.ServiceMeta) error {
	serviceToClientMap.Delete(meta.ServiceName)
	return nil
}

func GetClient(serviceName string) (*myClient, bool) {
	return serviceToClientMap.Get(serviceName)
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
	case "roundRobin":
		return client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()), nil
	default:
		return client.Option{}, fmt.Errorf("load balance choice %s not supported", lbChoice)
	}
}
