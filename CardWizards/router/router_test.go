package router

import (
	"net/http"
	"testing"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/common"
	"github.com/yiwen101/CardWizards/common/descriptor"
)

/* tested: GetRouteManager, isAnnotatedRoute, Getroute */
func TestValidateAnnotatedRoutes(t *testing.T) {
	// initialise descriptor manager, which this module depend on
	descriptor.BuildDescriptorManager(common.RelativePathToIDLFromTest2)
	dm, err := descriptor.GetDescriptorManager()
	test.Assert(t, err == nil)

	// test GetRouteManager
	rmTemp, err := GetRouteManager()
	test.Assert(t, err == nil)
	test.Assert(t, rmTemp != nil)

	// service name is arithmetic, method name is Add, but annotated path is /arith/add
	//Response Add(1: Request request ) ( api.get = "/arith/add" )
	rm := routeManagerImpl{dm: dm}

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

	routes, err := rmTemp.BuildAndGetRoute()
	test.Assert(t, err == nil)
	test.Assert(t, len(routes) > 0)
}
