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

func TestGetDescriptorFromPath(t *testing.T) {
	filename := "arithmatic.thrift"
	includeDir := RelativePathToIDL
	d, e := buildDescriptorKeeperFromPath(filename, includeDir)
	test.Assert(t, e == nil)

	des := d.get()
	test.Assert(t, des != nil)

	fuc, e := des.LookupFunctionByMethod("Add")
	test.Assert(t, e == nil)
	test.Assert(t, fuc != nil)

	fuc, e = des.LookupFunctionByMethod("fake")
	test.Assert(t, e != nil)
	test.Assert(t, fuc == nil)
}

/* struct
-> type: struct?
  -> struct - for entry in fieldByIDlen {
	if optional, next
	else if name not included, false
	else recall check type
  }
  -> i64
*/
