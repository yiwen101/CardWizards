package configuer

import (
	"os"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	loadbalance "github.com/cloudwego/kitex/pkg/loadbalance"

	"github.com/kitex-contrib/registry-nacos/resolver"
)

type serviceInfo struct {
	serviceName string
	IdlName    string
	IdlPath	string
	client 	genericclient.Client
	provider generic.DescriptorProvider
	

}

type idlsServicesManagerImplement struct {
}
type idlsServicesManager interface{
}


var serviceToClientMap = make(map[string]genericclient.Client)

func buildGenericClients() {
	serviceToClientMap = make(map[string]genericclient.Client)

	thiriftFiles, err := os.ReadDir("../IDL/")
	if err != nil {
		hlog.Fatal("failure reading thrrift files at IDL directory: %v", err)
	}

	nacosResolver, err := resolver.NewDefaultNacosResolver()
	if err != nil {
		hlog.Fatalf("err in building nacos resolver, please check your nacos server is on:%v", err)
	}

	for _, file := range thiriftFiles {
		if file.IsDir() {
			hlog.Fatal("failure reading thrrift files at IDL directory as it contains directory")
		}

		if file.Name()[len(file.Name())-7:] != ".thrift" {
			hlog.Fatal("failure reading thrrift files at IDL directory as it contains non-thrift file %s", file.Name())
		}

		// Get service name by deleting ".thrift" from the end of the file name
		serviceName := file.Name()[:len(file.Name())-7]

		p, err := generic.NewThriftFileProvider(file.Name(), "../../IDL/")
		if err != nil {
			hlog.Fatalf("new thrift provider failed: %v", err)
			break
		}
		g, err := generic.HTTPThriftGeneric(p)
		if err != nil {
			hlog.Fatal(err)
		}

		client, err := genericclient.NewClient(
			serviceName,
			g,
			client.WithResolver(nacosResolver),
			client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()),
		)
		if err != nil {
			hlog.Fatal("error in building generic client for service %s: %v", serviceName, err)
		}

		serviceToClientMap[serviceName] = client
	}
	hlog.Info("generic clients built successfully")
}
