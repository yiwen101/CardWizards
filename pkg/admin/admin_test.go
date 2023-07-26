package admin

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/hertz-contrib/cors"
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

	bytes, code = call("GET", "/admin/route/arithmetic/Add")
	test.Assert(t, code == http.StatusOK)
	var rInfo []routeInfo
	err = sonic.Unmarshal(bytes, &rInfo)
	test.Assert(t, err == nil, err)

	bytes, code = call("GET", "/admin/proxy")
	test.Assert(t, code == http.StatusOK)
	var b bool
	err = sonic.Unmarshal(bytes, &b)
	test.Assert(t, err == nil, err)

	bytes, code = call("GET", "/admin/lb/arithmetic")
	test.Assert(t, code == http.StatusOK)
	var s string
	err = sonic.Unmarshal(bytes, &s)
	test.Assert(t, err == nil, err)
	_, code = call("PUT", "/admin/lb/arithmetic/random")
	test.Assert(t, code == http.StatusOK)
	serviceMeta, err := store.InfoStore.GetServiceInfo("arithmetic")
	test.Assert(t, err == nil)
	test.Assert(t, serviceMeta.LbType == "random")

	_, code = callWithBody("PUT", "/admin/proxy", false)
	test.Assert(t, code == http.StatusOK)
	isOn, err := store.InfoStore.CheckProxyStatus()
	test.Assert(t, !isOn)
	test.Assert(t, err == nil, err)

	_, code = call("DELETE", "/admin/service/arithmetic")
	test.Assert(t, code == http.StatusOK)
	_, err = store.InfoStore.GetServiceInfo("arithmetic")
	test.Assert(t, err != nil)
	_, code = call("POST", "/admin/service/arithmetic.thrift/cluster1")
	test.Assert(t, code == http.StatusOK)
	_, err = store.InfoStore.GetServiceInfo("arithmetic")
	test.Assert(t, err == nil)

	_, code = callWithBody("PUT", "/admin/api/arithmetic/Add", false)
	test.Assert(t, code == http.StatusOK)
	meta, err := store.InfoStore.CheckAPIStatus("arithmetic", "Add")
	test.Assert(t, err == nil)
	test.Assert(t, !meta.IsOn)

	_, code = callWithBody("PUT", "/admin/service/arithmetic", true)
	test.Assert(t, code == http.StatusOK)
	meta, err = store.InfoStore.CheckAPIStatus("arithmetic", "Add")
	test.Assert(t, err == nil)
	test.Assert(t, meta.IsOn)

	callWithBody("PUT", "/admin/validation/arithmetic/Add", true)
	meta, err = store.InfoStore.CheckAPIStatus("arithmetic", "Add")
	test.Assert(t, err == nil)
	test.Assert(t, meta.ValidationOn)

	_, code = call("PUT", "/admin/service/arithmetic/arithmetic.thrift/mockCluster")
	test.Assert(t, code == http.StatusOK)
	serviceMeta, err = store.InfoStore.GetServiceInfo("arithmetic")
	test.Assert(t, err == nil)
	test.Assert(t, serviceMeta.ClusterName == "mockCluster")
	_, code = call("POST", "/admin/service/fake.thrift/cluster1")
	test.Assert(t, code != 200)

	_, code = callWithBody("POST", "/admin/route/arithmetic/Add", routeInfo{HttpMethod: "GET", Url: "test"})
	test.Assert(t, code == http.StatusOK)
	_, ok := router.GetRoute("GET", "/test")
	test.Assert(t, ok)

	_, code = callWithBody("PUT", "/admin/route/arithmetic/Add", routeUpdate{routeInfo{HttpMethod: "GET", Url: "test"}, routeInfo{HttpMethod: "GET", Url: "test2"}})
	test.Assert(t, code == http.StatusOK)
	_, ok = router.GetRoute("GET", "/test2")
	test.Assert(t, ok)
	_, ok = router.GetRoute("GET", "/test")
	test.Assert(t, !ok)

	_, code = callWithBody("DELETE", "/admin/route/arithmetic/Add", routeInfo{HttpMethod: "GET", Url: "test2"})
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
	admin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},                   // Update this to match your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // Add the allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Add the allowed request headers
		ExposeHeaders:    []string{"Content-Length"},                          // Expose additional response headers if needed
		AllowCredentials: true,                                                // Allow credentials (e.g., cookies, authorization headers)
		MaxAge:           12 * time.Hour,                                      // Set the preflight request cache duration
	}))
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
				c.String(http.StatusInternalServerError, err.Error())
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
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, translateService(s))
		})
	admin.GET("/api/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.CheckAPIStatus(c.Param("serviceName"), c.Param("methodName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			c.JSON(http.StatusOK, translateAPI(s))
		})
	admin.GET("/api/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.GetServiceInfo(c.Param("serviceName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			apis := []apiInfo{}
			for _, api := range s.APIs {
				apis = append(apis, translateAPI(api))
			}

			c.JSON(http.StatusOK, apis)
		})
	admin.GET("/route/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.CheckAPIStatus(c.Param("serviceName"), c.Param("methodName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			c.JSON(http.StatusOK, buildRoutesFromApi(s))
		})
	admin.GET("/proxy",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.CheckProxyStatus()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, s)
		})
	admin.GET("/lb/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.GetLbType(c.Param("serviceName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, translateLB(s))
		})
	admin.PUT("/proxy",
		func(ctx context.Context, c *app.RequestContext) {
			bytes, err := c.Body()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			var b bool
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			if b {
				err = store.InfoStore.TurnOnProxy()
			} else {
				err = store.InfoStore.TurnOffProxy()
			}
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Proxy status updated")
		})
	admin.PUT("/service/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			bytes, err := c.Body()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			var b bool
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			if b {
				err = store.InfoStore.TurnOnService(c.Param("serviceName"))
			} else {
				err = store.InfoStore.TurnOffService(c.Param("serviceName"))
			}
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			c.String(http.StatusOK, "Service status updated")
		})
	admin.PUT("/api/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			bytes, err := c.Body()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			var b bool
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			if b {
				err = store.InfoStore.TurnOnAPI(c.Param("serviceName"), c.Param("methodName"))
			} else {
				err = store.InfoStore.TurnOffAPI(c.Param("serviceName"), c.Param("methodName"))
			}
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			c.String(http.StatusOK, "Service status updated")
		})
	admin.PUT("/validation/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			bytes, err := c.Body()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			var b bool
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			if b {
				err = store.InfoStore.TurnOnValidation(c.Param("serviceName"), c.Param("methodName"))
			} else {
				err = store.InfoStore.TurnOffValidation(c.Param("serviceName"), c.Param("methodName"))
			}
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			c.String(http.StatusOK, "Validation status updated")
		})
	admin.POST("/service/:idlFileName/:clusterName",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.AddService(c.Param("idlFileName"), c.Param("clusterName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Service is added")
		})
	admin.PUT("/service/:serviceName/:idlFileName/:clusterName",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.UpdateService(c.Param("serviceName"), c.Param("idlFileName"), c.Param("clusterName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Service is updated")
		})
	admin.DELETE("/service/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.RemoveService(c.Param("serviceName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Service is removed")
		})
	admin.POST("/route/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			bytes, err := c.Body()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			var b routeInfo
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			err = store.InfoStore.AddRoute(c.Param("serviceName"), c.Param("methodName"), b.HttpMethod, "/"+b.Url)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Route is added")
		})
	admin.PUT("/route/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			bytes, err := c.Body()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			var b routeUpdate
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			err = store.InfoStore.ModifyRoute(c.Param("serviceName"), c.Param("methodName"), b.OldRoute.HttpMethod, "/"+b.OldRoute.Url, b.NewRoute.HttpMethod, "/"+b.NewRoute.Url)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Route is updated")
		})
	admin.DELETE("/route/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			bytes, err := c.Body()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			var b routeInfo
			err = sonic.Unmarshal(bytes, &b)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			err = store.InfoStore.RemoveRoute(c.Param("serviceName"), c.Param("methodName"), b.HttpMethod, "/"+b.Url)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Route is removed")
		})
	admin.PUT("/lb/:serviceName/:lbType",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.SetLbType(c.Param("serviceName"), c.Param("lbType"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "LbType is set")
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
