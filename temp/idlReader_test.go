package temp

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
)

func TestBuildGenericClientFromPath(t *testing.T) {
	filename := "arithmatic.thrift"
	includeDir := RelativePathToIDL
	_, e := buildGenericClientFromPath(filename, includeDir)
	test.Assert(t, e == nil)
}

// thrift package name = service name

func TestBuildGenericClients(t *testing.T) {
	err := buildGenericClients()
	test.Assert(t, err == nil)
	g1 := serviceToClientMap["arithmatic"]
	test.Assert(t, g1 != nil)
	g2 := serviceToClientMap["fake"]
	test.Assert(t, g2 == nil)
}
