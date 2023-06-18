package client

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/common"
)

func TestBuildGenericClientFromPath(t *testing.T) {
	filename := "arithmatic.thrift"
	includeDir := common.RelativePathToIDL
	_, e := buildGenericClientFromPath(filename, includeDir)
	test.Assert(t, e == nil)
}

// thrift package name = service name

func TestBuildGenericClients(t *testing.T) {
	err := BuildGenericClients()
	test.Assert(t, err == nil)
	g1 := ServiceToClientMap["arithmatic"]
	test.Assert(t, g1 != nil)
	g2 := ServiceToClientMap["fake"]
	test.Assert(t, g2 == nil)
}
