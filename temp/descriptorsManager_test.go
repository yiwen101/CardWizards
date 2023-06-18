package temp

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
)

func TestDescriptorsManager(t *testing.T) {
	error := buildDescriptorManager()
	test.Assert(t, error == nil)
	error = DescsManager.ValidateServiceAndMethodName("arithmatic", "Add")
	test.Assert(t, error == nil)
	error = DescsManager.ValidateServiceAndMethodName("arithmatic", "fake")
	test.Assert(t, error != nil)
}
