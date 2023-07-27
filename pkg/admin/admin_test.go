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
)

var testRouter *route.Engine

// todo adjust according to the change in the frontend
func TestAdmin(t *testing.T) {
	store.InfoStore.Load("proxyAddress", "../../testing/idl", "")
	testRouter = route.NewEngine(config.NewOptions([]config.Option{}))
	testRegisterAdmin(testRouter)
	testRegisterProxy(testRouter)

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
	store.InfoStore.Load("proxyAddress", "../../testing/idl", "password")
	testRouter = route.NewEngine(config.NewOptions([]config.Option{}))
	testRegisterAdmin(testRouter)
	testRegisterProxy(testRouter)

	w := ut.PerformRequest(
		testRouter,
		"GET",
		"/admin/proxy",
		&ut.Body{},
		ut.Header{Key: "Content-Type", Value: "application/json"},
	)
	code := w.Result().StatusCode()
	test.Assert(t, code != 200)

	w = ut.PerformRequest(
		testRouter,
		"GET",
		"/admin/proxy",
		&ut.Body{},
		ut.Header{Key: "Content-Type", Value: "application/json"},
		ut.Header{Key: "Password", Value: "password"},
	)
	code = w.Result().StatusCode()
	test.Assert(t, code == 200)

}

// make sure this function's body is indentical to the function Register
func testRegisterAdmin(r *route.Engine) {
	admin := r.Group("/admin")
	/*
		admin.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			//AllowOrigins:     []string{"http://localhost:3000"},                   // Update this to match your frontend URL
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // Add the allowed HTTP methods
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Add the allowed request headers
			ExposeHeaders:    []string{"Content-Length"},                          // Expose additional response headers if needed
			AllowCredentials: true,                                                // Allow credentials (e.g., cookies, authorization headers)
			MaxAge:           12 * time.Hour,                                      // Set the preflight request cache duration
		}))
	*/
	password := store.InfoStore.Password
	if password != "" {
		admin.Use(func(ctx context.Context, c *app.RequestContext) {
			if string(c.GetHeader("Password")) != password {
				c.AbortWithMsg("wrong password", http.StatusBadRequest)
			}
		})
	}

	admin.GET("/service",
		func(ctx context.Context, c *app.RequestContext) {
			services, err := store.InfoStore.GetAllServiceNames()
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			result := []serviceInfo{}

			for _, meta := range services {
				result = append(result, translateService(meta))
			}
			c.JSON(http.StatusOK, result)
		})
	admin.GET("/service/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.GetServiceInfo(c.Param("serviceName"))
			if err != nil {
				c.JSON(http.StatusNotFound, err.Error())
				return
			}
			c.JSON(http.StatusOK, translateService(s))
		})
	admin.POST("/service",
		func(ctx context.Context, c *app.RequestContext) {
			bytes, err := c.Body()
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			var b serviceInfo
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			serviceName, err := store.InfoStore.AddService(b.IdlFileName, b.ClusterName)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			serviceMeta, err := store.InfoStore.GetServiceInfo(serviceName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			newLb, ok := parseLB(b.LoadBalanceOption)
			if !ok {
				c.JSON(http.StatusBadRequest, "invalid load balance option")
				return
			}
			if newLb != serviceMeta.LbType {
				err = store.InfoStore.SetLbType(serviceName, newLb)
				if err != nil {
					c.JSON(http.StatusBadRequest, err.Error())
					return
				}
			}
			if b.IsSleeping {
				err = store.InfoStore.TurnOffService(serviceName)
				if err != nil {
					c.JSON(http.StatusBadRequest, err.Error())
					return
				}
			}
			serviceMeta, _ = store.InfoStore.GetServiceInfo(serviceName)

			c.JSON(http.StatusOK, translateService(serviceMeta))
		})
	admin.PUT("/service/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {

			bytes, err := c.Body()
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			var b serviceInfo
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			if b.ServiceName != c.Param("serviceName") {
				c.JSON(http.StatusBadRequest, "service name does not match")
				return
			}
			serviceMeta, err := store.InfoStore.GetServiceInfo(c.Param("serviceName"))
			if err != nil {
				c.JSON(http.StatusBadRequest, "invalid service name")
				return
			}

			newLb, ok := parseLB(b.LoadBalanceOption)
			if !ok {
				c.JSON(http.StatusBadRequest, "invalid load balance option")
				return
			}
			ServiceName := b.ServiceName
			if b.ClusterName != serviceMeta.ClusterName || b.ReloadIDL {
				ServiceName, err = store.InfoStore.UpdateService(b.ServiceName, b.IdlFileName, b.ClusterName)
				if err != nil {
					c.JSON(http.StatusBadRequest, err.Error())
					return
				}
				serviceMeta, err = store.InfoStore.GetServiceInfo(ServiceName)
				if err != nil {
					c.JSON(http.StatusInternalServerError, err.Error())
					return
				}
			}

			if newLb != serviceMeta.LbType {
				err = store.InfoStore.SetLbType(ServiceName, newLb)
				if err != nil {
					c.JSON(http.StatusInternalServerError, err.Error())
					return
				}
			}

			if b.IsSleeping != isSleeping(serviceMeta) {
				if !b.IsSleeping {
					err = store.InfoStore.TurnOnService(c.Param("serviceName"))
				} else {
					err = store.InfoStore.TurnOffService(c.Param("serviceName"))
				}
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			serviceMeta, err = store.InfoStore.GetServiceInfo(ServiceName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			c.JSON(http.StatusOK, translateService(serviceMeta))

		})
	admin.DELETE("/service/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.RemoveService(c.Param("serviceName"))
			if err != nil {
				c.JSON(http.StatusNotFound, err.Error())
				return
			}
			c.JSON(http.StatusOK, "Service is removed")
		})

	admin.GET("/api/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.CheckAPIStatus(c.Param("serviceName"), c.Param("methodName"))
			if err != nil {
				c.JSON(http.StatusNotFound, err.Error())
				return
			}

			c.JSON(http.StatusOK, translateAPI(s))
		})

	admin.GET("/api/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.GetServiceInfo(c.Param("serviceName"))
			if err != nil {
				c.JSON(http.StatusNotFound, err.Error())
				return
			}
			apis := []apiInfo{}
			for _, api := range s.APIs {
				apis = append(apis, translateAPI(api))
			}

			c.JSON(http.StatusOK, apis)
		})
	admin.GET("/api",
		func(ctx context.Context, c *app.RequestContext) {
			services, err := store.InfoStore.GetAllServiceNames()
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			apis := []apiInfo{}
			for _, service := range services {
				for _, api := range service.APIs {
					apis = append(apis, translateAPI(api))
				}
			}
			c.JSON(http.StatusOK, apis)
		})

	admin.PUT("/api/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			apiMeta, err := store.InfoStore.CheckAPIStatus(c.Param("serviceName"), c.Param("methodName"))
			if err != nil {
				c.JSON(http.StatusNotFound, err.Error())
				return
			}

			bytes, err := c.Body()
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			var b apiInfo
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			if b.ServiceName != c.Param("serviceName") || b.MethodName != c.Param("methodName") {
				c.JSON(http.StatusBadRequest, "service name or method name does not match")
				return
			}

			if !b.IsSleeping != apiMeta.IsOn {
				if !b.IsSleeping {
					err = store.InfoStore.TurnOnAPI(c.Param("serviceName"), c.Param("methodName"))
				} else {
					err = store.InfoStore.TurnOffAPI(c.Param("serviceName"), c.Param("methodName"))
				}
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			if b.ValidationStatus != apiMeta.ValidationOn {
				if b.ValidationStatus {
					err = store.InfoStore.TurnOnValidation(c.Param("serviceName"), c.Param("methodName"))
				} else {
					err = store.InfoStore.TurnOffValidation(c.Param("serviceName"), c.Param("methodName"))
				}
			}

			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			newApi, err := store.InfoStore.CheckAPIStatus(c.Param("serviceName"), c.Param("methodName"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			c.JSON(http.StatusOK, translateAPI(newApi))
		})

	admin.GET("/proxy",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.CheckProxyStatus()
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, s)
		})
	admin.PUT("/proxy",
		func(ctx context.Context, c *app.RequestContext) {
			bytes, err := c.Body()
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			var b bool
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			if b {
				err = store.InfoStore.TurnOnProxy()
			} else {
				err = store.InfoStore.TurnOffProxy()
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, "Proxy status updated")
		})

	admin.GET("/route",
		func(ctx context.Context, c *app.RequestContext) {
			services, err := store.InfoStore.GetAllServiceNames()
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			routes := []routeInfo{}
			for _, service := range services {
				routes = append(routes, buildRoutesfromService(service)...)
			}
			c.JSON(http.StatusOK, routes)
		})
	admin.GET("/route/:httpMethod/*url",
		func(ctx context.Context, c *app.RequestContext) {
			method, url := c.Param("httpMethod"), "/"+c.Param("url")

			r, ok := router.GetRoute(method, url)
			if !ok {
				c.JSON(http.StatusNotFound, "route not found")
				return
			}
			c.JSON(http.StatusOK, routeInfo{method + url, r.ServiceName, r.MethodName, method, url})
		},
	)
	admin.POST("/route",
		func(ctx context.Context, c *app.RequestContext) {
			bytes, err := c.Body()
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			var b routeInfo
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}

			err = store.InfoStore.AddRoute(b.ServiceName, b.MethodName, b.HttpMethod, b.Url)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			b.Id = b.HttpMethod + b.Url
			c.JSON(http.StatusOK, b)
		})
	admin.PUT("/route/:httpMethod/*url",
		func(ctx context.Context, c *app.RequestContext) {
			method, url := c.Param("httpMethod"), "/"+c.Param("url")

			r, ok := router.GetRoute(method, url)
			if !ok {
				c.JSON(http.StatusNotFound, "route not found")
				return
			}
			bytes, err := c.Body()
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			var b routeInfo
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			if b.MethodName != r.MethodName || b.ServiceName != r.ServiceName {
				c.JSON(http.StatusBadRequest, "route info does not match")
				return
			}
			err = store.InfoStore.ModifyRoute(b.ServiceName, b.MethodName, method, url, b.HttpMethod, b.Url)
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			c.JSON(http.StatusOK, b)
		})
	admin.DELETE("/route/:httpMethod/*url",
		func(ctx context.Context, c *app.RequestContext) {
			method, url := c.Param("httpMethod"), "/"+c.Param("url")
			r, ok := router.GetRoute(method, url)
			if !ok {
				c.JSON(http.StatusNotFound, "route not found")
				return
			}

			err := store.InfoStore.RemoveRoute(r.ServiceName, r.MethodName, method, url)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, "Route is removed")
		})
}

func testRegisterProxy(r *route.Engine) {
	r.GET("/*:test",
		func(ctx context.Context, c *app.RequestContext) {
			proxy.Proxy.Serve(ctx, c, nil)
		})
}

func call(method string, url string) ([]byte, int) {
	w := ut.PerformRequest(testRouter, method, url, &ut.Body{}, ut.Header{Key: "Content-Type", Value: "application/json"})
	return w.Result().Body(), w.Result().StatusCode()
}

func callWithBody(method string, url string, body interface{}) ([]byte, int) {
	bs, _ := sonic.Marshal(body)
	b := bytes.NewBuffer(bs)
	w := ut.PerformRequest(testRouter, method, url, &ut.Body{Body: b, Len: b.Len()}, ut.Header{Key: "Content-Type", Value: "application/json"})
	return w.Result().Body(), w.Result().StatusCode()
}
