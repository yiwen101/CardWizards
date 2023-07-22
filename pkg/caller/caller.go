package caller

/*
This package is about offering the right generic client. This class should only be responsible for keeping clients
up to date, for example when update the api by uploading a new idl file, or updating the load balance choice
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
Choose to make my provider as the provider in the generic package seems to block infinitely after emiting
only one instance of descriptor. But building a generic client require a channel of descriptor. The easy
solution of rereading and parsing the file in not only expensive, but also violate the single source of truth
principle. The service descriptors in the store package should be the only source of truth, and serviceMeta,
apiMeta, and the whole bunch of other data (for example, client) are generared with them. So we implement
our own provider to build the generic client
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

An alternative is to include the handler logic here so to reduce passing parameters via funciton on
copy. However, it will mess up with the single responsibility principle. This class should only be
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
	case "roundrobin":
		return client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()), nil
	default:
		return client.Option{}, fmt.Errorf("load balance choice %s not supported", lbChoice)
	}
}
