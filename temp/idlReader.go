package temp

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/kitex-contrib/registry-nacos/resolver"
)

type descriptorKeeper struct {
	svcDsc   atomic.Value
	provider generic.DescriptorProvider
	codec    remote.PayloadCodec
}

func newDescriptorKeeper(p generic.DescriptorProvider) (*descriptorKeeper, error) {
	svc := <-p.Provide()
	d := &descriptorKeeper{provider: p}
	d.svcDsc.Store(svc)
	go d.update()
	return d, nil
}
func (d *descriptorKeeper) update() {
	for {
		svc, ok := <-d.provider.Provide()
		if !ok {
			return
		}
		d.svcDsc.Store(svc)
	}
}
func (d *descriptorKeeper) get() *descriptor.ServiceDescriptor {
	svcDsc, ok := d.svcDsc.Load().(*descriptor.ServiceDescriptor)
	if !ok {
		hlog.Fatalf("invalid service descriptor")
	}
	return svcDsc
}

var container map[string]*descriptorKeeper

var serviceToClientMap map[string]genericclient.Client

func getDescriptorFromPath(fileName, includeDir string) (*descriptorKeeper, error) {
	p, err := generic.NewThriftFileProvider(fileName, includeDir)
	if err != nil {
		hlog.Fatalf("new thrift provider failed: %v", err)
	}
	descriptor, err := newDescriptorKeeper(p)
	if err != nil {
		hlog.Fatalf("new descriptor keeper failed: %v", err)
	}
	return descriptor, err
}

func buildContainer() error {
	container = make(map[string]*descriptorKeeper)

	thiriftFiles, err := os.ReadDir(RelativePathToIDL)
	if err != nil {
		hlog.Fatal("failure reading thrrift files at IDL directory: %v", err)
	}
	flag := false

	for _, file := range thiriftFiles {
		if file.IsDir() {
			hlog.Fatal("failure reading thrrift files at IDL directory as it contains directory")
		}

		if file.Name()[len(file.Name())-7:] != ".thrift" {
			hlog.Fatal("failure reading thrrift files at IDL directory as it contains non-thrift file %s", file.Name())
		}

		d, err := getDescriptorFromPath(file.Name(), "../IDL/")
		if err != nil {
			flag = true
			hlog.Fatal("error in building descriptor for service %s: %v", file.Name(), err)
		}
		serviceName := file.Name()[:len(file.Name())-7]
		container[serviceName] = d
	}

	if flag {
		return fmt.Errorf("error in building generic clients")
	} else {
		hlog.Info("generic container built successfully")
		return nil
	}
}

func buildGenericClientFromPath(fileName, includeDir string, opts ...client.Option) (genericclient.Client, error) {
	//serviceToClientMap = make(map[string]genericclient.Client)

	p, err := generic.NewThriftFileProvider(fileName, includeDir)
	if err != nil {
		hlog.Fatalf("new thrift provider failed: %v", err)
	}

	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		hlog.Fatalf("new JSONThriftGeneric failed: %v", err)
	}

	client, err := genericclient.NewClient(
		fileName,
		g,
		opts...,
	)
	if err != nil {
		hlog.Fatal("error in building generic client for service %s: %v", fileName, err)
	}

	return client, err
}

func buildGenericClients() error {
	serviceToClientMap = make(map[string]genericclient.Client)

	thiriftFiles, err := os.ReadDir(RelativePathToIDL)
	if err != nil {
		hlog.Fatal("failure reading thrrift files at IDL directory: %v", err)
	}

	nacosResolver, err := resolver.NewDefaultNacosResolver()
	if err != nil {
		hlog.Fatalf("err in building nacos resolver, please check your nacos server is on:%v", err)
	}

	flag := false

	for _, file := range thiriftFiles {
		if file.IsDir() {
			hlog.Fatal("failure reading thrrift files at IDL directory as it contains directory")
		}

		if file.Name()[len(file.Name())-7:] != ".thrift" {
			hlog.Fatal("failure reading thrrift files at IDL directory as it contains non-thrift file %s", file.Name())
		}

		client, err := buildGenericClientFromPath(
			file.Name(),
			"../IDL/",
			client.WithResolver(nacosResolver),
			client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()),
		)

		// Get service name by deleting ".thrift" from the end of the file name
		serviceName := file.Name()[:len(file.Name())-7]

		if err != nil {
			flag = true
			hlog.Fatal("error in building generic client for service %s: %v", serviceName, err)
		}
		serviceToClientMap[serviceName] = client
	}
	if flag {
		return fmt.Errorf("error in building generic clients")
	} else {
		hlog.Info("generic clients built successfully")
		return nil
	}
}

func validate(serviceName, mathodName string) bool {
	_, ok := container[serviceName]
	return ok
}
