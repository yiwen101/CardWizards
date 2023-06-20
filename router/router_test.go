package router

import (
	"net/http"
	"testing"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/common/descriptor"
)

func TestValidateAnnotatedRoutes(t *testing.T) {
	// initialise descriptor manager, which this module depend on
	descriptor.BuildDescriptorManager("../IDL")
	dm, err := descriptor.GetDescriptorManager()
	test.Assert(t, err == nil)

	// test GetRouteManager
	rmTemp, err := GetRouteManager()
	test.Assert(t, err == nil)
	test.Assert(t, rmTemp != nil)

	// service name is arithmatic, method name is Add, but annotated path is /arith/add
	//Response Add(1: Request request ) ( api.get = "/arith/add" )
	rm := routeManagerImpl{dm: dm}

	httpReq, err := http.NewRequest("GET", "/arithmatic/Add", nil)
	test.Assert(t, err == nil)
	req, err := generic.FromHTTPRequest(httpReq)
	test.Assert(t, err == nil)
	serviceName, methodName, err := rm.isAnnotatedRoute(req)
	test.Assert(t, err != nil)
	test.Assert(t, serviceName == "" && methodName == "")

	httpReq, err = http.NewRequest("GET", "/arith/add", nil)
	test.Assert(t, err == nil)
	req, err = generic.FromHTTPRequest(httpReq)
	test.Assert(t, err == nil)
	serviceName, methodName, err = rm.isAnnotatedRoute(req)
	test.Assert(t, err == nil)
	test.Assert(t, serviceName == "arithmatic" && methodName == "Add")

}

func TestGetRoute(t *testing.T) {
	descriptor.BuildDescriptorManager("../IDL")

	rmTemp, err := GetRouteManager()
	test.Assert(t, err == nil)
	test.Assert(t, rmTemp != nil)

	routes, err := rmTemp.GetRoutes()
	test.Assert(t, err == nil)
	test.Assert(t, len(routes) > 0)
}
