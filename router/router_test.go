package router

import (
	"net/http"
	"testing"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/common"
	"github.com/yiwen101/CardWizards/pkg/store/descriptor"
)

/* tested: GetRouteManager, isAnnotatedRoute, Getroute */
func TestValidateAnnotatedRoutes(t *testing.T) {
	// initialise descriptor manager, which this module depend on
	descriptor.BuildDescriptorManager(common.RelativePathToIDLFromTest2)
	//dm, err := descriptor.GetDescriptorManager()
	//test.Assert(t, err == nil)

	// test GetRouteManager
	rmTemp, err := GetRouteManager()
	rmTemp.InitRoute()
	test.Assert(t, err == nil)
	test.Assert(t, rmTemp != nil)

	// service name is arithmetic, method name is Add, but annotated path is /arith/add
	//Response Add(1: Request request ) ( api.get = "/arith/add" )
	rm := rmTemp.(*routeManagerImpl)

	httpReq, err := http.NewRequest("GET", "/arithmetic/Add", nil)
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
	test.Assert(t, serviceName == "arithmetic" && methodName == "Add")
}

func TestBuildAndGetRoutes(t *testing.T) {
	descriptor.BuildDescriptorManager(common.RelativePathToIDLFromTest2)

	rmTemp, err := GetRouteManager()
	test.Assert(t, err == nil)
	test.Assert(t, rmTemp != nil)

	err = rmTemp.InitRoute()
	test.Assert(t, err == nil)
	_, ok := rmTemp.GetRoute("POST", "/arith/add")
	test.Assert(t, !ok)
	api, ok := rmTemp.GetRoute("GET", "/arithmetic/Add")
	test.Assert(t, ok)
	test.Assert(t, api.ServiceName == "arithmetic" && api.MethodName == "Add")
}

func TestUpdateAndDeleteRoute(t *testing.T) {
	descriptor.BuildDescriptorManager(common.RelativePathToIDLFromTest2)

	rmTemp, err := GetRouteManager()
	test.Assert(t, err == nil)
	test.Assert(t, rmTemp != nil)

	err = rmTemp.InitRoute()
	test.Assert(t, err == nil)
	_, ok := rmTemp.GetRoute("POST", "/arith/add")
	test.Assert(t, !ok)
	api, ok := rmTemp.GetRoute("GET", "/arithmetic/Add")
	test.Assert(t, ok)
	test.Assert(t, api.ServiceName == "arithmetic" && api.MethodName == "Add")
	rmTemp.UpdateRoute("GET", "/arithmetic/Add", "GET", "/arith/add")
	_, ok = rmTemp.GetRoute("GET", "/arithmetic/Add")
	test.Assert(t, !ok)
	_, ok = rmTemp.GetRoute("GET", "/arith/add")
	test.Assert(t, ok)
	rmTemp.DeleteRoute("GET", "/arith/add")
	_, ok = rmTemp.GetRoute("GET", "/arith/add")
	test.Assert(t, !ok)
}
