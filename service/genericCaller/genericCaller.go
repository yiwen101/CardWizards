package service

import (
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	"github.com/kitex-contrib/registry-nacos/resolver"
)

var ServiceNamesToGenericClients = make(map[string]genericclient.Client)

func initMap() {
	if ServiceNamesToGenericClients == nil {
		ServiceNamesToGenericClients = make(map[string]genericclient.Client)
	}
}

func ReadIDLsFromPath(relativePath string) {
	initMap()
	dirEntries, err := os.ReadDir(relativePath)
	if err != nil {
		hlog.Fatalf("new thrift file provider failed: %v", err)
	}

	for _, entry := range dirEntries {
		// Ignore directories and non-thrift files
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		if fileName[len(fileName)-5:] != ".thrift" {
			continue
		}

		// Get service name by deleting ".thrift" from the end of the file name
		serviceName := fileName[:len(fileName)-7]

		provider, err := generic.NewThriftFileProvider(fileName, relativePath)
		if err != nil {
			hlog.Fatalf("new thrift provider failed: %v", err)
			break
		}
		g, err := generic.HTTPThriftGeneric(provider)
		if err != nil {
			hlog.Fatal(err)
		}

		nacosResolver, err := resolver.NewDefaultNacosResolver()
		if err != nil {
			hlog.Fatalf("err:%v", err)
		}

		cli, err := genericclient.NewClient(
			serviceName,
			g,
			client.WithResolver(nacosResolver),
			client.WithTransportProtocol(transport.TTHeader),
			// what this line means?
			client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		)
		if err != nil {
			hlog.Fatal(err)
		}

		ServiceNamesToGenericClients[serviceName] = cli

	}
}
