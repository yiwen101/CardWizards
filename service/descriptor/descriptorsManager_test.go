package descriptor

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/common"
)

func TestDescriptorsManager(t *testing.T) {
	error := BuildDescriptorManager(common.RelativePathToIDLFromTest)
	test.Assert(t, error == nil)
	error = DescsManager.ValidateServiceAndMethodName("arithmatic", "Add")
	test.Assert(t, error == nil)
	error = DescsManager.ValidateServiceAndMethodName("arithmatic", "fake")
	test.Assert(t, error != nil)
}
