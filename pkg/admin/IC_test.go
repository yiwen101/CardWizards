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

func TestAdmin(t *testing.T) {
	store.InfoStore.Load("proxyAddress", "../../testing/idl", "")
	testRouter = route.NewEngine(config.NewOptions([]config.Option{}))
	TestregisterAdmin(testRouter)
	TestregisterProxy(testRouter)

	bytes, code := call("GET", "/admin/service")
	test.Assert(t, code == http.StatusOK)
	var j map[string]interface{}
	err := sonic.Unmarshal(bytes, &j)
	test.Assert(t, err == nil, err)

	bytes, code = call("GET", "/admin/service/arithmetic")
	test.Assert(t, code == http.StatusOK)
	j = make(map[string]interface{})
	err = sonic.Unmarshal(bytes, &j)
	test.Assert(t, err == nil, err)

	bytes, code = call("GET", "/admin/api/arithmetic/Add")
	test.Assert(t, code == http.StatusOK)
	j = make(map[string]interface{})
	err = sonic.Unmarshal(bytes, &j)
	test.Assert(t, err == nil, err)

	bytes, code = call("GET", "/admin/api/arithmetic")
	test.Assert(t, code == http.StatusOK)
	j = make(map[string]interface{})
	err = sonic.Unmarshal(bytes, &j)
	test.Assert(t, err == nil, err)

	bytes, code = call("GET", "/admin/route/arithmetic/Add")
	test.Assert(t, code == http.StatusOK)
	a := []struct{ httpMethod, url string }{}
	err = sonic.Unmarshal(bytes, &a)
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

	_, code = call("POST", "/admin/route/arithmetic/Add/GET/test")
	test.Assert(t, code == http.StatusOK)
	_, ok := router.GetRoute("GET", "/test")
	test.Assert(t, ok)

	_, code = call("PUT", "/admin/route/arithmetic/Add/GET/test/GET/test2")
	test.Assert(t, code == http.StatusOK)
	_, ok = router.GetRoute("GET", "/test2")
	test.Assert(t, ok)
	_, ok = router.GetRoute("GET", "/test")
	test.Assert(t, !ok)

	_, code = call("DELETE", "/admin/route/arithmetic/Add/GET/test2")
	test.Assert(t, code == http.StatusOK)
	_, ok = router.GetRoute("GET", "/test2")
	test.Assert(t, !ok)
}

func TestPassword(t *testing.T) {
	store.InfoStore.Load("proxyAddress", "../../testing/idl", "password")
	testRouter = route.NewEngine(config.NewOptions([]config.Option{}))
	TestregisterAdmin(testRouter)
	TestregisterProxy(testRouter)

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

func TestregisterAdmin(r *route.Engine) {
	admin := r.Group("/admin")

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

			m, err := store.InfoStore.GetAllServiceNames()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, m)
		})
	admin.GET("/service/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.GetServiceInfo(c.Param("serviceName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, s)
		})
	admin.GET("/api/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.CheckAPIStatus(c.Param("serviceName"), c.Param("methodName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			routes := []struct {
				httpMethod string
				url        string
			}{}
			for httpMethod, value := range s.Routes {
				for url := range value {
					routes = append(routes, struct {
						httpMethod string
						url        string
					}{httpMethod, url})
				}
			}

			result := struct {
				serviceName string
				methodName  string
				apiStatus   bool
				routes      []struct {
					httpMethod string
					url        string
				}
				validationStatus bool
			}{s.ServiceName, s.MethodName, s.IsOn, routes, s.ValidationOn}
			c.JSON(http.StatusOK, result)
		})
	admin.GET("/api/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.GetAPIs(c.Param("serviceName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, s)
		})
	admin.GET("/route/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.GetRoutes(c.Param("serviceName"), c.Param("methodName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			result := []struct {
				httpMethod string
				url        string
			}{}
			for httpMethod, value := range s {
				for url := range value {
					result = append(result, struct {
						httpMethod string
						url        string
					}{httpMethod, url})
				}
			}

			c.JSON(http.StatusOK, result)
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
			c.JSON(http.StatusOK, s)
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
	admin.POST("/route/:serviceName/:methodName/:httpMethod/:url",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.AddRoute(c.Param("serviceName"), c.Param("methodName"), c.Param("httpMethod"), "/"+c.Param("url"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Route is added")
		})
	admin.PUT("/route/:serviceName/:methodName/:httpMethod/:url/:newHttpMethod/:newUrl",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.ModifyRoute(c.Param("serviceName"), c.Param("methodName"), c.Param("httpMethod"), "/"+c.Param("url"), c.Param("newHttpMethod"), "/"+c.Param("newUrl"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Route is updated")
		})
	admin.DELETE("/route/:serviceName/:methodName/:httpMethod/:url",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.RemoveRoute(c.Param("serviceName"), c.Param("methodName"), c.Param("httpMethod"), "/"+c.Param("url"))
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

func TestregisterProxy(r *route.Engine) {
	r.GET("/*:test",
		func(ctx context.Context, c *app.RequestContext) {
			proxy.Proxy.Serve(ctx, c, nil)
		})
}

func call(method string, url string) ([]byte, int) {
	w := ut.PerformRequest(testRouter, method, url, &ut.Body{}, ut.Header{Key: "Content-Type", Value: "application/json"})
	return w.Result().Body(), w.Result().StatusCode()
}

func callWithBody(method string, url string, body bool) ([]byte, int) {
	bs, _ := sonic.Marshal(body)
	b := bytes.NewBuffer(bs)
	w := ut.PerformRequest(testRouter, method, url, &ut.Body{Body: b, Len: b.Len()}, ut.Header{Key: "Content-Type", Value: "application/json"})
	return w.Result().Body(), w.Result().StatusCode()
}
