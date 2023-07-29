package admin

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/pkg/proxy"
	"github.com/yiwen101/CardWizards/pkg/router"
	"github.com/yiwen101/CardWizards/pkg/store"
	"github.com/yiwen101/CardWizards/pkg/utils"
)

var testRouter *route.Engine
var admin *route.RouterGroup
var defaultHeaders = []ut.Header{}

func init() {
	store.InfoStore.Load("proxyAddress", utils.PkgToIDL, "password")
	testRouter = route.NewEngine(config.NewOptions([]config.Option{}))
	registerAPIGateway(testRouter)

	admin = testRouter.Group("/admin")
	for _, f := range registerlist {
		f(admin)
	}
	defaultHeaders = []ut.Header{
		{Key: "Content-Type", Value: "application/json"},
		{Key: "Password", Value: store.InfoStore.Password},
	}
}

// turn on the kitex server and nacos server before run this test
func TestAdmin(t *testing.T) {
	bytes, code := call("GET", "/admin/service")
	test.Assert(t, code == http.StatusOK)
	var servs []serviceInfo
	err := sonic.Unmarshal(bytes, &servs)
	test.Assert(t, err == nil, err)

	bytes, code = call("GET", "/admin/service/arithmetic")
	test.Assert(t, code == http.StatusOK)
	var serv serviceInfo
	err = sonic.Unmarshal(bytes, &serv)
	test.Assert(t, err == nil, err)

	bytes, code = call("GET", "/admin/api/arithmetic/Add")
	test.Assert(t, code == http.StatusOK)
	var aInfo apiInfo
	err = sonic.Unmarshal(bytes, &aInfo)
	test.Assert(t, err == nil, err)

	bytes, code = call("GET", "/admin/api/arithmetic")
	test.Assert(t, code == http.StatusOK)
	var aInfos []apiInfo
	err = sonic.Unmarshal(bytes, &aInfos)
	test.Assert(t, err == nil, err)

	bytes, code = call("GET", "/admin/route/GET/arithmetic/Add")
	test.Assert(t, code == http.StatusOK)
	var rInfo routeInfo
	err = sonic.Unmarshal(bytes, &rInfo)
	test.Assert(t, err == nil, err)

	bytes, code = call("GET", "/admin/proxy")
	test.Assert(t, code == http.StatusOK)
	var b bool
	err = sonic.Unmarshal(bytes, &b)
	test.Assert(t, err == nil, err)

	_, code = callWithBody("PUT", "/admin/proxy", false)
	test.Assert(t, code == http.StatusOK)
	isOn, err := store.InfoStore.CheckProxyStatus()
	test.Assert(t, !isOn)
	test.Assert(t, err == nil, err)

	_, code = call("DELETE", "/admin/service/arithmetic")
	test.Assert(t, code == http.StatusOK)
	_, err = store.InfoStore.GetServiceInfo("arithmetic")
	test.Assert(t, err != nil)
	_, code = callWithBody("POST", "/admin/service", serviceInfo{ReloadIDL: true, IdlFileName: "arithmetic.thrift", ClusterName: "arithmetic", LoadBalanceOption: "weighted random"})
	test.Assert(t, code == http.StatusOK)
	_, err = store.InfoStore.GetServiceInfo("arithmetic")
	test.Assert(t, err == nil)

	_, code = callWithBody("PUT", "/admin/api/arithmetic/Add", apiInfo{ServiceName: "arithmetic", MethodName: "Add", IsSleeping: true, ValidationStatus: true})
	test.Assert(t, code == http.StatusOK)
	meta, err := store.InfoStore.CheckAPIStatus("arithmetic", "Add")
	test.Assert(t, err == nil)
	test.Assert(t, !meta.IsOn)
	test.Assert(t, meta.ValidationOn)

	_, code = callWithBody("PUT", "/admin/service/arithmetic", serviceInfo{ReloadIDL: true, ServiceName: "arithmetic", IdlFileName: "arithmetic.thrift", ClusterName: "mockCluster", LoadBalanceOption: "weighted random"})
	test.Assert(t, code == http.StatusOK)
	serviceMeta, err := store.InfoStore.GetServiceInfo("arithmetic")
	test.Assert(t, err == nil)
	test.Assert(t, serviceMeta.ClusterName == "mockCluster")
	_, code = callWithBody("POST", "/admin/service", serviceInfo{ReloadIDL: true, ServiceName: "arithmetic", IdlFileName: "arithmetic.protobuf", ClusterName: "arithmetic", LoadBalanceOption: "weighted random"})
	test.Assert(t, code != 200)

	_, code = callWithBody("POST", "/admin/route", routeInfo{ServiceName: "arithmetic", MethodName: "Add", HttpMethod: "GET", Url: "/test"})
	test.Assert(t, code == http.StatusOK)
	_, ok := router.GetRoute("GET", "/test")
	test.Assert(t, ok)

	_, code = callWithBody("PUT", "/admin/route/GET/test", routeInfo{ServiceName: "arithmetic", MethodName: "Add", HttpMethod: "GET", Url: "/test2"})
	test.Assert(t, code == http.StatusOK)
	_, ok = router.GetRoute("GET", "/test2")
	test.Assert(t, ok)
	_, ok = router.GetRoute("GET", "/test")
	test.Assert(t, !ok)

	_, code = call("DELETE", "/admin/route/GET/test2")
	test.Assert(t, code == http.StatusOK)
	_, ok = router.GetRoute("GET", "/test2")
	test.Assert(t, !ok)
}

func TestPassword(t *testing.T) {
	password := store.InfoStore.Password
	test.Assert(t, password != "")
	admin.Use(func(ctx context.Context, c *app.RequestContext) {
		if string(c.GetHeader("Password")) != password {
			c.AbortWithMsg("wrong password", http.StatusBadRequest)
		}
	})
	_, code := call("GET", "/admin/proxy")
	test.Assert(t, code == 200)

	_, code = callWithoutPassword("PUT", "/admin/proxy")
	test.Assert(t, code != 200)
}

func call(method string, url string) ([]byte, int) {
	w := ut.PerformRequest(testRouter, method, url, &ut.Body{}, defaultHeaders...)
	return w.Result().Body(), w.Result().StatusCode()
}

func callWithBody(method string, url string, body interface{}, headers ...ut.Header) ([]byte, int) {
	bs, _ := sonic.Marshal(body)
	b := bytes.NewBuffer(bs)
	headers = append(headers, defaultHeaders...)
	w := ut.PerformRequest(testRouter, method, url, &ut.Body{Body: b, Len: b.Len()}, headers...)
	return w.Result().Body(), w.Result().StatusCode()
}

func callWithoutPassword(method string, url string) ([]byte, int) {
	w := ut.PerformRequest(testRouter, method, url, &ut.Body{}, ut.Header{Key: "Content-Type", Value: "application/json"})
	return w.Result().Body(), w.Result().StatusCode()
}

func registerAPIGateway(r *route.Engine) {
	r.GET("/*:test",
		func(ctx context.Context, c *app.RequestContext) {
			proxy.Proxy.Serve(ctx, c, nil)
		})
}
