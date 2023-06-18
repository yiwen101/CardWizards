package descriptor

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
)

func TestDescriptorsManager(t *testing.T) {
	error := BuildDescriptorManager()
	test.Assert(t, error == nil)
	error = DescsManager.ValidateServiceAndMethodName("arithmatic", "Add")
	test.Assert(t, error == nil)
	error = DescsManager.ValidateServiceAndMethodName("arithmatic", "fake")
	test.Assert(t, error != nil)
}
