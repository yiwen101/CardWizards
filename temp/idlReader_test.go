package temp

import (
	"testing"

	"github.com/cloudwego/kitex/pkg/generic/descriptor"
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
	d, e := getDescriptorFromPath(filename, includeDir)
	des := d.get()
	fuc, e := des.LookupFunctionByMethod("Add")
	if e != nil {
		// means method not found
	}
	var ls map[int32]*descriptor.FieldDescriptor
	ls = fuc.Request.Struct.FieldsByID
	// find the number of keys for the map
	if ls[2] != nil {
		// invalid number of fields
	}
	req := ls[1]

	println(req)
	test.Assert(t, e == nil)
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
