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
}

func (d *descriptorsManagerImpl) GetServiceName(filename string) (string, error) {
	keeper, ok := d.m[filename]
	if !ok {
		return "", fmt.Errorf("fileName %s not found", filename)
	}
	return keeper.get().Name, nil
}

func (d *descriptorsManagerImpl) GetRouters() map[string]descriptor.Router {
	if d.routers == nil {
		d.buildRouters()
	}
	return d.routers
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
		log.Printf("reading file : %s", file.Name())
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
		descManager.m[file.Name()] = d
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
	// map of file names to descriptor keepers
	m map[string]*descriptorKeeper
	// map of service name to file name
	serivceMap map[string]string
	// map of service name to router
	routers map[string]descriptor.Router
}

func newDescriptorsManagerImpl() *descriptorsManagerImpl {
	return &descriptorsManagerImpl{m: make(map[string]*descriptorKeeper)}
}
func (d *descriptorsManagerImpl) buildServiceMap() {
	d.serivceMap = make(map[string]string)
	for fileName, manager := range d.m {
		d.serivceMap[manager.get().Name] = fileName
	}
}
func (d *descriptorsManagerImpl) getFileName(serviceName string) (string, error) {
	if d.serivceMap == nil {
		d.buildServiceMap()
	}

	fileName, ok := d.serivceMap[serviceName]
	if !ok {
		return "", fmt.Errorf("service %s not found", serviceName)
	}
	return fileName, nil
}

func (d *descriptorsManagerImpl) buildRouters() error {
	d.routers = make(map[string]descriptor.Router)
	for fileName, descriptorKeeper := range d.m {
		serviceName, err := d.getFileName(fileName)
		if err != nil {
			return err
		}
		d.routers[serviceName] = descriptorKeeper.get().Router
	}
	return nil
}
