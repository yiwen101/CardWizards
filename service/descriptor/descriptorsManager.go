package descriptor

import (
	"fmt"
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/yiwen101/CardWizards/common"
)

var DescsManager DescriptorsManager

type DescriptorsManager interface {
	ValidateServiceAndMethodNameWithAnnotedRoutes(req *descriptor.HTTPRequest) (string, error)
	ValidateServiceAndMethodName(serviceName, methodName string) error
	GetFunctionDescriptor(serviceName, methodName string) (*descriptor.FunctionDescriptor, error)
}

type descriptorsManagerImpl struct {
	m map[string]*descriptorKeeper
}

func newDescriptorsManagerImpl() *descriptorsManagerImpl {
	return &descriptorsManagerImpl{m: make(map[string]*descriptorKeeper)}
}

func (d *descriptorsManagerImpl) ValidateServiceAndMethodName(serviceName, methodName string) error {

	manager, ok := d.m[serviceName]
	if !ok {
		return fmt.Errorf("service %s not found", serviceName)
	}
	return manager.validateMethodName(methodName)
}

func (d *descriptorsManagerImpl) ValidateServiceAndMethodNameWithAnnotedRoutes(req *descriptor.HTTPRequest) (string, error) {
	for serviceName, manager := range d.m {
		if manager.matchedRouter(req) {
			return serviceName, nil
		}
	}
	return "", fmt.Errorf("service not found")
}

func (d *descriptorsManagerImpl) GetFunctionDescriptor(serviceName, methodName string) (*descriptor.FunctionDescriptor, error) {
	manager, ok := d.m[serviceName]
	if !ok {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}
	return manager.get().LookupFunctionByMethod(methodName)
}
func BuildDescriptorManager() error {
	descManager := newDescriptorsManagerImpl()
	thiriftFiles, err := os.ReadDir(common.RelativePathToIDL)
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

		d, err := buildDescriptorKeeperFromPath(file.Name(), "../IDL/")
		if err != nil {
			flag = true
			hlog.Fatal("error in building descriptor for service %s: %v", file.Name(), err)
		}
		serviceName := file.Name()[:len(file.Name())-7]
		descManager.m[serviceName] = d
	}

	DescsManager = descManager

	if flag {
		return fmt.Errorf("error in building generic clients")
	} else {
		hlog.Info("generic container built successfully")
		return nil
	}
}

