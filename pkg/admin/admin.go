package admin

import (
	"context"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yiwen101/CardWizards/pkg/router"
	"github.com/yiwen101/CardWizards/pkg/store"
)

func Register(r *server.Hertz) {
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

func translateLB(lb string) string {
	switch lb {
	case "default":
		return "weighted round robin"
	case "random":
		return "weighted random"
	case "roundRobin":
		return "weighted round robin"
	default:
		return "error load balance option"
	}
}

func parseLB(lb string) (string, bool) {
	switch lb {
	case "weighted round robin":
		return "default", true
	case "weighted random":
		return "random", true
	default:
		return "", false
	}
}

func buildRoutesFromApi(api *store.ApiMeta) []routeInfo {
	routes := []routeInfo{}

	for httpMethod, value := range api.Routes {
		for url := range value {
			routes = append(routes, routeInfo{httpMethod + url, api.ServiceName, api.MethodName, httpMethod, url})
		}
	}
	return routes
}
func buildRoutesfromService(mata *store.ServiceMeta) []routeInfo {
	routes := []routeInfo{}
	for _, api := range mata.APIs {
		routes = append(routes, buildRoutesFromApi(api)...)
	}
	return routes
}

func translateAPI(api *store.ApiMeta) apiInfo {

	return apiInfo{api.ServiceName + "/" + api.MethodName, api.ServiceName, api.MethodName, !api.IsOn, api.ValidationOn}
}

func translateService(s *store.ServiceMeta) serviceInfo {
	idlName, _ := store.InfoStore.GetIdlFileName(s.ServiceName)
	return serviceInfo{s.ServiceName, s.ServiceName, translateLB(s.LbType), s.ClusterName, isSleeping(s), idlName, false}
}

func isSleeping(s *store.ServiceMeta) bool {
	for _, api := range s.APIs {
		if api.IsOn {
			return false
		}
	}
	return true
}

type routeUpdate struct {
	OldRoute routeInfo `json:"oldRoute"`
	NewRoute routeInfo `json:"newRoute"`
}

type serviceInfo struct {
	Id                string `json:"id"`
	ServiceName       string `json:"serviceName"`
	LoadBalanceOption string `json:"loadBalanceOption"`
	ClusterName       string `json:"clusterName"`
	IsSleeping        bool   `json:"isSleeping"`
	IdlFileName       string `json:"idlFileName"`
	ReloadIDL         bool   `json:"reloadIdl"`
}

type routeInfo struct {
	Id          string `json:"id"`
	ServiceName string `json:"serviceName"`
	MethodName  string `json:"methodName"`
	HttpMethod  string `json:"httpMethod"`
	Url         string `json:"url"`
}

type apiInfo struct {
	Id               string `json:"id"`
	ServiceName      string `json:"serviceName"`
	MethodName       string `json:"methodName"`
	IsSleeping       bool   `json:"isSleeping"`
	ValidationStatus bool   `json:"validationOn"`
}
