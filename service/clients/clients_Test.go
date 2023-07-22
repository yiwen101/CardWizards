package clients

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/pkg/store/descriptor"
)

func TestBuildGenericClientFromPath(t *testing.T) {
	filename := "arithmetic.thrift"
	includeDir := "../../IDL"
	_, e := buildGenericClientFromPath("arithmetic", filename, includeDir)
	test.Assert(t, e == nil)
}

func TestBuildGenericClientsAndGetGenericClientforService(t *testing.T) {
	descriptor.BuildDescriptorManager("../../IDL")
	err := BuildGenericClients("../../IDL")
	test.Assert(t, err == nil)
	g1, err := ClientManager.GetClient("arithmetic")
	test.Assert(t, g1 != nil)
	test.Assert(t, err == nil)
	g2, err := ClientManager.GetClient("arithmetic2")
	test.Assert(t, g2 == nil)
	test.Assert(t, err != nil)
}

func TestProvideTwo(t *testing.T) {
	s1, s2, err := provideTwo()
	test.Assert(t, err == nil)
	test.Assert(t, s1 != nil)
	test.Assert(t, s2 != nil)
}
