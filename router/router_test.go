package router

import (
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	kitexDescriptor "github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/common/descriptor"
)

func TestRouter(t *testing.T) {
	// todo, how to do unit testing
	r := NewRouteManager()
	h := server.Default(
		server.WithHostPorts("127.0.0.1:8080"),
	)
	r.RegisterRoutes(h)
}

func Test(t *testing.T) {

	err := descriptor.BuildDescriptorManager("../IDL")
	test.Assert(t, err == nil)
	router, err := descriptor.DescriptorManager.GetRouterForService("arithmatic")
	test.Assert(t, err == nil)
	// service name is arithmatic, method name is Add, but annotated path is /arith/add
	//Response Add(1: Request request ) ( api.get = "/arith/add" )

	httpRequest := kitexDescriptor.HTTPRequest{}
	httpRequest.Method = "GET"
	httpRequest.Path = "/arithmatic/Add"
	funcDescriptor, err := router.Lookup(&httpRequest)
	test.Assert(t, err != nil)
	test.Assert(t, funcDescriptor == nil)
	httpRequest.Path = "/arith/add"
	funcDescriptor, err = router.Lookup(&httpRequest)
	test.Assert(t, err == nil)
	test.Assert(t, funcDescriptor != nil)
}
