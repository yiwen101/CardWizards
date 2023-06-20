package clients

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/common"
)

func TestBuildGenericClientFromPath(t *testing.T) {
	filename := "arithmatic.thrift"
	includeDir := common.RelativePathToIDLFromTest
	_, e := buildGenericClientFromPath(filename, includeDir)
	test.Assert(t, e == nil)
}

func TestBuildGenericClients(t *testing.T) {
	err := BuildGenericClients(common.RelativePathToIDLFromTest)
	test.Assert(t, err == nil)
	g1, err := GetGenericClientforService("arithmatic")
	test.Assert(t, g1 != nil)
	test.Assert(t, err == nil)
	g2, err := GetGenericClientforService("arithmatic2")
	test.Assert(t, g2 == nil)
	test.Assert(t, err != nil)
}
