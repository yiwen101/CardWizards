package store

import (
	"fmt"
	"sync/atomic"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	//"github.com/cloudwego/kitex/pkg/remote"
)

type descriptorKeeper struct {
	fileName string
	svcDsc   atomic.Value
	provider generic.DescriptorProvider
	//codec       remote.PayloadCodec
}

func newDescriptorKeeper(p generic.DescriptorProvider, filename string) (*descriptorKeeper, error) {
	svc := <-p.Provide()
	d := &descriptorKeeper{provider: p, fileName: filename}
	d.svcDsc.Store(svc)
	go d.update()
	return d, nil
}
func (d *descriptorKeeper) update() {
	for {
		svc, ok := <-d.provider.Provide()
		if !ok {
			return
		}
		d.svcDsc.Store(svc)
	}
}
func (d *descriptorKeeper) Get() (*descriptor.ServiceDescriptor, error) {
	svcDsc, ok := d.svcDsc.Load().(*descriptor.ServiceDescriptor)
	if !ok {
		return nil, fmt.Errorf("invalid service descriptor for %s", d.fileName)
	}
	return svcDsc, nil
}
func (d *descriptorKeeper) GetFileName() (string, error) {
	return d.fileName, nil
}
func (d *descriptorKeeper) validateMethodName(methodName string) error {
	sd, err := d.Get()
	if err != nil {
		return err
	}
	_, err = sd.LookupFunctionByMethod(methodName)
	if err != nil {
		return fmt.Errorf("method %s not found", methodName)
	}
	return nil
}

func buildDescriptorKeeperFromPath(fileName, includeDir string) (*descriptorKeeper, error) {

	p, err := generic.NewThriftFileProvider(fileName, includeDir)
	if err != nil {
		return nil, err
	}
	descriptor, err := newDescriptorKeeper(p, fileName)
	if err != nil {
		return nil, err
	}
	return descriptor, err
}
