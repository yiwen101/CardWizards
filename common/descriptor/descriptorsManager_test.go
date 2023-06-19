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
