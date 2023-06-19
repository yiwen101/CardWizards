package router

import (
	"testing"

	kitexDescriptor "github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/common/descriptor"
)

func Test(t *testing.T) {
	// relative path from here to IDL
	descriptor.BuildDescriptorManager("../IDL")
	dm, err := descriptor.GetDescriptorManager()
	test.Assert(t, err == nil)

	// service name is arithmatic, method name is Add, but annotated path is /arith/add
	//Response Add(1: Request request ) ( api.get = "/arith/add" )

	httpRequest := kitexDescriptor.HTTPRequest{}
	httpRequest.Method = "GET"
	httpRequest.Path = "/arithmatic/Add"
	serviceName, methodName, err := dm.GetMathchedRouterName(&httpRequest)
	test.Assert(t, err != nil)
	test.Assert(t, serviceName == "")
	test.Assert(t, methodName == "")
	httpRequest.Path = "/arith/add"
	serviceName, methodName, err = dm.GetMathchedRouterName(&httpRequest)
	test.Assert(t, err == nil)
	test.Assert(t, serviceName == "arithmatic")
	test.Assert(t, methodName == "Add")
}
