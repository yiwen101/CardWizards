package router

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/pkg/store"
	"github.com/yiwen101/CardWizards/pkg/utils"
)

func TestBuildGetUpdateAndDeleteRoutes(t *testing.T) {
	store.InfoStore.Load("", utils.PkgToIDL, "")
	data, ok := GetRoute("GET", "/arithmetic/Add")
	test.Assert(t, ok)
	test.Assert(t, data.ServiceName == "arithmetic" && data.MethodName == "Add")
	store.InfoStore.RemoveRoute("arithmetic", "Add", "GET", "/arithmetic/Add")
	data, ok = GetRoute("GET", "/arithmetic/Add")
	test.Assert(t, !ok)
	store.InfoStore.AddRoute("arithmetic", "Add", "GET", "/arithmetic/Add")
	data, ok = GetRoute("GET", "/arithmetic/Add")
	test.Assert(t, ok)
	test.Assert(t, data.ServiceName == "arithmetic" && data.MethodName == "Add")
	store.InfoStore.ModifyRoute("arithmetic", "Add", "GET", "/arithmetic/Add", "GET", "/test")
	data, ok = GetRoute("GET", "/arithmetic/Add")
	test.Assert(t, !ok)
	data, ok = GetRoute("GET", "/test")
	test.Assert(t, ok)
	test.Assert(t, data.ServiceName == "arithmetic" && data.MethodName == "Add")
}

/* test for route generated from thrift annotation. No longer in use. Maybe will fix in the future
func TestValidateAnnotatedRoutes(t *testing.T) {
	// initialise descriptor manager, which this module depend on
	descriptor.BuildDescriptorManager(utils.RelativePathToIDLFromTest2)
	//dm, err := descriptor.GetDescriptorManager()
	//test.Assert(t, err == nil)

	// test GetRouteManager

	// service name is arithmetic, method name is Add, but annotated path is /arith/add
	//Response Add(1: Request request ) ( api.get = "/arith/add" )

	httpReq, err := http.NewRequest("GET", "/arithmetic/Add", nil)
	test.Assert(t, err == nil)
	req, err := generic.FromHTTPRequest(httpReq)
	test.Assert(t, err == nil)
	serviceName, methodName, err := isAnnotatedRoute(req)
	test.Assert(t, err != nil)
	test.Assert(t, serviceName == "" && methodName == "")

	httpReq, err = http.NewRequest("GET", "/arith/add", nil)
	test.Assert(t, err == nil)
	req, err = generic.FromHTTPRequest(httpReq)
	test.Assert(t, err == nil)
	serviceName, methodName, err = isAnnotatedRoute(req)
	test.Assert(t, err == nil)
	test.Assert(t, serviceName == "arithmetic" && methodName == "Add")
}
*/
