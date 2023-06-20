package descriptor

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
)

type DescsManager interface {
	GetFunctionDescriptor(serviceName, methodName string) (*descriptor.FunctionDescriptor, error)
	GetServiceDescriptor(serviceName string) (*descriptor.ServiceDescriptor, error)
	GetServiceName(filename string) (string, error)
	GetRouters() map[string]descriptor.Router
	GetAllServiceNames() ([]string, error)
	GetAllMethodNames(serviceName string) ([]string, error)
}

func BuildDescriptorManager(relativePath string) error {
	descManager := newDescriptorsManagerImpl()
	thiriftFiles, err := os.ReadDir(relativePath)
	if err != nil {
		return err
	}

	for _, file := range thiriftFiles {
		log.Printf("reading file : %s", file.Name())
		if file.IsDir() {
			return fmt.Errorf("failure reading thrrift files at IDL directory as it contains directory")
		}

		if file.Name()[len(file.Name())-7:] != ".thrift" {
			return fmt.Errorf("failure reading thrrift files at IDL directory as it contains non-thrift file %s", file.Name())
		}

		d, err := buildDescriptorKeeperFromPath(file.Name(), relativePath)
		if err != nil {
			return fmt.Errorf("error in building descriptor for service %s: %v", file.Name(), err)
		}
		descManager.m[file.Name()] = d
	}

	descriptorManager = descManager

	hlog.Info("generic container built successfully")
	return nil

}

func GetDescriptorManager() (DescsManager, error) {
	if descriptorManager == nil {
		return nil, fmt.Errorf("descriptor manager not built")
	}
	return descriptorManager, nil
}

func (d *descriptorsManagerImpl) GetServiceName(filename string) (string, error) {
	keeper, ok := d.m[filename]
	if !ok {
		return "", fmt.Errorf("fileName %s not found", filename)
	}

	funcDesc, err := keeper.get()
	if err != nil {
		return "", err
	}

	return funcDesc.Name, nil
}

func (d *descriptorsManagerImpl) GetRouters() map[string]descriptor.Router {
	if d.routers == nil {
		d.buildRouters()
	}
	return d.routers
}

func (d *descriptorsManagerImpl) GetFunctionDescriptor(serviceName, methodName string) (*descriptor.FunctionDescriptor, error) {
	ser, err := d.GetServiceDescriptor(serviceName)
	if err != nil {
		return nil, err
	}
	return ser.LookupFunctionByMethod(methodName)
}

func (d *descriptorsManagerImpl) GetServiceDescriptor(serviceName string) (*descriptor.ServiceDescriptor, error) {
	str, err := d.getFileName(serviceName)
	if err != nil {
		return nil, err
	}
	descriptorKeeper, ok := d.m[str]
	if !ok {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}
	return descriptorKeeper.get()
}

func (d *descriptorsManagerImpl) GetAllServiceNames() ([]string, error) {
	result := make([]string, 0)
	for fileName := range d.m {
		serviceName, err := d.GetServiceName(fileName)
		if err != nil {
			return nil, err
		}
		result = append(result, serviceName)
	}
	return result, nil
}

func (d *descriptorsManagerImpl) GetAllMethodNames(serviceName string) ([]string, error) {
	serviceDescriptor, err := d.GetServiceDescriptor(serviceName)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)

	for methodName := range serviceDescriptor.Functions {
		result = append(result, methodName)
	}
	return result, nil
}
