package clients

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/kitex-contrib/registry-nacos/resolver"
	desc "github.com/yiwen101/CardWizards/common/descriptor"
	"github.com/yiwen101/CardWizards/configuer/clientOption"
)

var dm desc.DescsManager

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
		return err
	}

	for _, file := range thiriftFiles {
		if file.IsDir() {
			return fmt.Errorf("failure reading thrrift files at IDL directory as it contains directory")
		}

		if file.Name()[len(file.Name())-7:] != ".thrift" {
			return fmt.Errorf("failure reading thrrift files at IDL directory as it contains non-thrift file %s", file.Name())
		}

		// Get service name by deleting ".thrift" from the end of the file name
		serviceName, err := dm.GetServiceName(file.Name())
		if err != nil {
			return err
		}

		registryOption, err := getServiceRegistryOption(serviceName)
		if err != nil {
			return err
		}

		client, err := buildGenericClientFromPath(
			serviceName,
			file.Name(),
			relativePath,
			registryOption,
			clientOption.GetServiceLoadBalancerOption(serviceName),
		)

		if err != nil {
			return err
		}
		serviceToClientMapTemp[serviceName] = client
	}

	serviceToClientMap = serviceToClientMapTemp
	ClientManager = clientManager{
		mu:                 sync.RWMutex{},
		serviceToClientMap: serviceToClientMap,
	}

	return nil
}

var serviceToClientMap map[string]genericclient.Client

type clientManager struct {
	mu                 sync.RWMutex
	serviceToClientMap map[string]genericclient.Client
}

func (cm *clientManager) GetClient(serviceName string) (genericclient.Client, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	cli, ok := cm.serviceToClientMap[serviceName]
	if !ok {
		return nil, fmt.Errorf("no client found for service %s", serviceName)
	}
	return cli, nil
}

func (cm *clientManager) UpdateClient(serviceName, fileName, includeDir string, opts ...client.Option) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cli, err := buildGenericClientFromPath(serviceName, fileName, includeDir, opts...)
	if err != nil {
		return err
	}
	cm.serviceToClientMap[serviceName] = cli
	return nil
}

func buildGenericClientFromPath(serviceName, fileName, includeDir string, opts ...client.Option) (genericclient.Client, error) {
	pwd, _ := os.Getwd()

	p, err := generic.NewThriftFileProvider(fileName, includeDir)
	if err != nil {
		return nil, fmt.Errorf("failure reading thrift files at IDL directory: %s, pwd is %s, parameters are %s, %s", err.Error(), pwd, fileName, includeDir)
	}

	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		return nil, err
	}

	client, err := genericclient.NewClient(
		serviceName,
		g,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	log.Printf("client for service %s built, parameters are %s, %s\n", serviceName, fileName, includeDir)

	return client, err
}

func getServiceRegistryOption(serviceName string) (client.Option, error) {
	nacosResolver, err := resolver.NewDefaultNacosResolver()
	if err != nil {
		return client.Option{}, err
	}

	return client.WithResolver(nacosResolver), nil
}

var ClientManager clientManager
