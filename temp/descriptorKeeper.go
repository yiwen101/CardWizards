package temp

import (
	"fmt"
	"sync/atomic"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/cloudwego/kitex/pkg/remote"
)

type descriptorKeeper struct {
	svcDsc   atomic.Value
	provider generic.DescriptorProvider
	codec    remote.PayloadCodec
}

func newDescriptorKeeper(p generic.DescriptorProvider) (*descriptorKeeper, error) {
	svc := <-p.Provide()
	d := &descriptorKeeper{provider: p}
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
func (d *descriptorKeeper) get() *descriptor.ServiceDescriptor {
	svcDsc, ok := d.svcDsc.Load().(*descriptor.ServiceDescriptor)
	if !ok {
		hlog.Fatalf("invalid service descriptor")
	}
	return svcDsc
}
func (d *descriptorKeeper) validateMethodName(methodName string) error {
	_, err := d.get().LookupFunctionByMethod(methodName)
	if err != nil {
		return fmt.Errorf("method %s not found", methodName)
	}
	return nil
}

func buildDescriptorKeeperFromPath(fileName, includeDir string) (*descriptorKeeper, error) {
	p, err := generic.NewThriftFileProvider(fileName, includeDir)
	if err != nil {
		hlog.Fatalf("new thrift provider failed: %v", err)
	}
	descriptor, err := newDescriptorKeeper(p)
	if err != nil {
		hlog.Fatalf("new descriptor keeper failed: %v", err)
	}
	return descriptor, err
}
