package descriptor

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
)

type DescsManager interface {
	GetMathchedRouterName(req *descriptor.HTTPRequest) (string, string, error)
	GetFunctionDescriptor(serviceName, methodName string) (*descriptor.FunctionDescriptor, error)
	GetServiceDescriptor(serviceName string) (*descriptor.ServiceDescriptor, error)
}

func (d *descriptorsManagerImpl) GetMathchedRouterName(req *descriptor.HTTPRequest) (string, string, error) {
	// cache the path -> service/method?
	for serviceName, manager := range d.m {
		if methodname, match := manager.matchedRouter(req); match {
			return serviceName, methodname, nil
		}
	}
	return "", "", fmt.Errorf("service not found")
}

func (d *descriptorsManagerImpl) GetFunctionDescriptor(serviceName, methodName string) (*descriptor.FunctionDescriptor, error) {
	manager, ok := d.m[serviceName]
	if !ok {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}
	return manager.get().LookupFunctionByMethod(methodName)
}

func (d *descriptorsManagerImpl) GetServiceDescriptor(serviceName string) (*descriptor.ServiceDescriptor, error) {
	descriptorKeeper, ok := d.m[serviceName]
	if !ok {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}
	return descriptorKeeper.get(), nil
}

func BuildDescriptorManager(relativePath string) error {
	descManager := newDescriptorsManagerImpl()
	thiriftFiles, err := os.ReadDir(relativePath)
	if err != nil {
		hlog.Fatal("failure reading thrift files at IDL directory: %v", err)
	}
	flag := false

	for _, file := range thiriftFiles {
		log.Printf("file name: %s", file.Name())
		if file.IsDir() {
			hlog.Fatal("failure reading thrrift files at IDL directory as it contains directory")
		}

		if file.Name()[len(file.Name())-7:] != ".thrift" {
			hlog.Fatal("failure reading thrrift files at IDL directory as it contains non-thrift file %s", file.Name())
		}

		d, err := buildDescriptorKeeperFromPath(file.Name(), relativePath)
		if err != nil {
			flag = true
			hlog.Fatal("error in building descriptor for service %s: %v", file.Name(), err)
		}
		serviceName := file.Name()[:len(file.Name())-7]
		descManager.m[serviceName] = d
	}

	for serviceName := range descManager.m {
		router := descManager.routers[serviceName]
		descManager.routers[serviceName] = router
	}

	descriptorManager = descManager

	if flag {
		return fmt.Errorf("error in building generic clients")
	} else {
		hlog.Info("generic container built successfully")
		return nil
	}

}

func GetDescriptorManager() (DescsManager, error) {
	if descriptorManager == nil {
		return nil, fmt.Errorf("descriptor manager not built")
	}
	return descriptorManager, nil
}

var descriptorManager DescsManager

type descriptorsManagerImpl struct {
	m       map[string]*descriptorKeeper
	routers map[string]descriptor.Router
}

func newDescriptorsManagerImpl() *descriptorsManagerImpl {
	return &descriptorsManagerImpl{m: make(map[string]*descriptorKeeper), routers: make(map[string]descriptor.Router)}
}
