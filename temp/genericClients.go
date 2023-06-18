package temp

import (
	"fmt"
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/kitex-contrib/registry-nacos/resolver"
)

var ServiceToClientMap map[string]genericclient.Client

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
	ServiceToClientMap = make(map[string]genericclient.Client)

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
		ServiceToClientMap[serviceName] = client
	}
	if flag {
		return fmt.Errorf("error in building generic clients")
	} else {
		hlog.Info("generic clients built successfully")
		return nil
	}
}
