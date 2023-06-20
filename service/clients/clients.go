package clients

import (
	"fmt"
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/kitex-contrib/registry-nacos/resolver"
	desc "github.com/yiwen101/CardWizards/common/descriptor"
)

var dm desc.DescsManager

func GetGenericClientforService(serviceName string) (genericclient.Client, error) {
	if client, ok := serviceToClientMap[serviceName]; ok {
		return client, nil
	} else {
		return nil, fmt.Errorf("no client found for service %s", serviceName)
	}
}

func BuildGenericClients(relativePath string) error {
	if serviceToClientMap != nil {
		return nil
	}

	if dm == nil {
		dmTemp, err := desc.GetDescriptorManager()
		if err != nil {
			return err
		}
		dm = dmTemp
	}

	serviceToClientMapTemp := make(map[string]genericclient.Client)
	thiriftFiles, err := os.ReadDir(relativePath)
	if err != nil {
		hlog.Fatal("failure reading thrrift files at IDL directory: %v", err)
	}

	for _, file := range thiriftFiles {
		if file.IsDir() {
			hlog.Fatal("failure reading thrrift files at IDL directory as it contains directory")
		}

		if file.Name()[len(file.Name())-7:] != ".thrift" {
			hlog.Fatal("failure reading thrrift files at IDL directory as it contains non-thrift file %s", file.Name())
		}

		// Get service name by deleting ".thrift" from the end of the file name
		serviceName, err := dm.GetServiceName(file.Name())
		if err != nil {
			return err
		}

		client, err := buildGenericClientFromPath(
			serviceName,
			file.Name(),
			relativePath,
			getServiceRegistryOption(serviceName),
			getServiceLoadBalancerOption(serviceName),
		)

		if err != nil {
			hlog.Fatal("error in building generic client for service %s: %v", serviceName, err)
		}
		serviceToClientMapTemp[serviceName] = client
	}

	serviceToClientMap = serviceToClientMapTemp

	return nil
}

var serviceToClientMap map[string]genericclient.Client

func buildGenericClientFromPath(serviceName, fileName, includeDir string, opts ...client.Option) (genericclient.Client, error) {

	p, err := generic.NewThriftFileProvider(fileName, includeDir)
	if err != nil {
		hlog.Fatalf("new thrift provider failed: %v", err)
	}

	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		hlog.Fatalf("new JSONThriftGeneric failed: %v", err)
	}

	client, err := genericclient.NewClient(
		serviceName,
		g,
		opts...,
	)
	if err != nil {
		hlog.Fatal("error in building generic client for service %s: %v", fileName, err)
	}

	return client, err
}

func getServiceRegistryOption(serviceName string) client.Option {
	nacosResolver, err := resolver.NewDefaultNacosResolver()
	if err != nil {
		hlog.Fatalf("err in building nacos resolver, please check your nacos server is on:%v", err)
	}

	return client.WithResolver(nacosResolver)
}

func getServiceLoadBalancerOption(serviceName string) client.Option {
	// todo: enable optioning
	return client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer())
}
