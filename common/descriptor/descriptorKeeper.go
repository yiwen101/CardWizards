package descriptor

import (
	"fmt"
	"sync/atomic"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	//"github.com/cloudwego/kitex/pkg/remote"
)

type descriptorKeeper struct {
	serviceName string
	svcDsc      atomic.Value
	provider    generic.DescriptorProvider
	//codec       remote.PayloadCodec
}

func newDescriptorKeeper(p generic.DescriptorProvider, name string) (*descriptorKeeper, error) {
	svc := <-p.Provide()
	d := &descriptorKeeper{provider: p, serviceName: name}
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

func (d *descriptorKeeper) matchedRouter(req *descriptor.HTTPRequest) (string, bool) {
	router := d.get().Router
	if router == nil {
		return "", false
	}
	des, err := router.Lookup(req)
	if err == nil {
		return des.Name, true
	}
	return "", false
}

func buildDescriptorKeeperFromPath(fileName, includeDir string) (*descriptorKeeper, error) {
	serviceName := fileName[:len(fileName)-7]

	p, err := generic.NewThriftFileProvider(fileName, includeDir)
	if err != nil {
		hlog.Fatalf("new thrift provider failed: %v", err)
	}
	descriptor, err := newDescriptorKeeper(p, serviceName)
	if err != nil {
		hlog.Fatalf("new descriptor keeper failed: %v", err)
	}
	return descriptor, err
}
