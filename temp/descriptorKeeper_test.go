package temp

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
)

func TestDescriptorKeeper(t *testing.T) {
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

/*

func experiment(t *testing.T) {
	filename := "arithmatic.thrift"
	includeDir := RelativePathToIDL
	d, e := buildDescriptorKeeperFromPath(filename, includeDir)
	test.Assert(t, e == nil)

	des := d.get()
	test.Assert(t, des != nil)
	node := des.Router.
}


	/*
Router
tree
GET prefix/ppath, function/Name

*/
