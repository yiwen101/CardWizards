package descriptor

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/common"
)

func TestDescriptorsManager(t *testing.T) {
	BuildDescriptorManager(common.RelativePathToIDLFromTest)
	dm, error := GetDescriptorManager()
	test.Assert(t, error == nil)
	_, error = dm.GetServiceDescriptor("arithmetic")
	test.Assert(t, error == nil)
	_, error = dm.GetFunctionDescriptor("arithmetic", "fake")
	test.Assert(t, error != nil)
	_, error = dm.GetFunctionDescriptor("arithmetic", "Add")
	test.Assert(t, error == nil)
}

func TestNilRouter(t *testing.T) {
	BuildDescriptorManager(common.RelativePathToIDLFromTest)
	dm, err := GetDescriptorManager()
	test.Assert(t, err == nil)
	service, err := dm.GetServiceDescriptor("arithmetic")
	test.Assert(t, err == nil)
	test.Assert(t, service != nil)
	serviceName := service.Name
	test.Assert(t, serviceName == "arithmetic")
	funcs := service.Functions
	test.Assert(t, len(funcs) == 5)
	for _, f := range funcs {
		str := f.Name
		test.Assert(t, str != "")
	}

}
