package caller

import (
	"fmt"
	"sync"

	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/yiwen101/CardWizards/pkg/store"
)

var Caller caller

func init() {
	Caller = caller{
		store:              store.InfoStore,
		mu:                 sync.RWMutex{},
		serviceToClientMap: make(map[string]genericclient.Client),
	}
	Caller.store.RegisterServiceMapListener(&serviceChangeHandler{})
}

type serviceChangeHandler struct{}

func (s *serviceChangeHandler) OnStatechanged(data ...interface{}) error {
	changeType := data[0].(string)
	serviceName := data[1].(string)
	if changeType == "add" {
		return Caller.AddClient(serviceName)
	}
	if changeType == "delete" {
		return Caller.DeleteClient(serviceName)
	}
	return fmt.Errorf("unknown change type %s", changeType)
}

type caller struct {
	store              *store.Store
	mu                 sync.RWMutex
	serviceToClientMap map[string]genericclient.Client
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
func (c *caller) AddClient(serviceName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	meta, err := c.store.GetServiceInfo(serviceName)
	if err != nil {
		return err
	}
	decriptor, err := meta.Descriptor.Get()
	if err != nil {
		return err
	}
	p, err := newMyProvider(decriptor)
	if err != nil {
		return fmt.Errorf("error makring myProvider for %s: %s", serviceName, err.Error())
	}

	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		return err
	}

	client, err := genericclient.NewClient(
		serviceName,
		g,
		//opts...,
	)
	if err != nil {
		return err
	}
	c.serviceToClientMap[serviceName] = client

	return nil
}

func (c *caller) DeleteClient(serviceName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.serviceToClientMap, serviceName)
	return nil
}

/*
func (c *caller) UpdateClient(serviceName string) error {
	err := c.DeleteClient(serviceName)
	if err != nil {
		return err
	}
	return c.AddClient(serviceName)
}
*/

func (cm *caller) GetClient(serviceName string) (genericclient.Client, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	cli, ok := cm.serviceToClientMap[serviceName]
	if !ok {
		return nil, nil
	}
	return cli, nil
}
