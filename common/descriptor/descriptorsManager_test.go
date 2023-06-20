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
	_, error = dm.GetServiceDescriptor("arithmatic")
	test.Assert(t, error == nil)
	_, error = dm.GetFunctionDescriptor("arithmatic", "fake")
	test.Assert(t, error != nil)
	_, error = dm.GetFunctionDescriptor("arithmatic", "Add")
	test.Assert(t, error == nil)
}

func TestNilRouter(t *testing.T) {
	BuildDescriptorManager(common.RelativePathToIDLFromTest)
	dm, err := GetDescriptorManager()
	test.Assert(t, err == nil)
	service, err := dm.GetServiceDescriptor("arithmatic")
	test.Assert(t, err == nil)
	test.Assert(t, service != nil)
	serviceName := service.Name
	test.Assert(t, serviceName == "arithmatic")
	funcs := service.Functions
	test.Assert(t, len(funcs) == 5)
	for _, f := range funcs {
		str := f.Name
		test.Assert(t, str != "")
	}

}
