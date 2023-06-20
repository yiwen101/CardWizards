package descriptor

import (
	"fmt"

	"github.com/cloudwego/kitex/pkg/generic/descriptor"
)

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
		serviceName, err := d.GetServiceName(fileName)
		if err != nil {
			return err
		}
		d.routers[serviceName] = descriptorKeeper.get().Router
	}
	return nil
}
