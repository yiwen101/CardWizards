package store

import (
	"fmt"
	"sync/atomic"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
)

// responsible for providing the service descriptor of a service.

type descriptorKeeper struct {
	fileName string
	svcDsc   atomic.Value
	provider generic.DescriptorProvider
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

func newDescriptorKeeper(p generic.DescriptorProvider, filename string) (*descriptorKeeper, error) {
	svc := <-p.Provide()
	d := &descriptorKeeper{provider: p, fileName: filename}
	d.svcDsc.Store(svc)
	go d.update()
	return d, nil
}

/*
it turned out the the provider in the generic package seems to block infinitely after emiting a single
instance of descriptor, so there is actually no need for running update. But we still keep this method
for future use.
*/
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
