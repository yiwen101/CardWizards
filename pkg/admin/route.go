package admin

import (
	"context"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/yiwen101/CardWizards/pkg/router"
	"github.com/yiwen101/CardWizards/pkg/store"
)

/*
List of APIs:
1. get all route names: GET /admin/route
2. get a route info: GET /admin/route/:routeID
3. add route: POST /admin/route                     json body: RouteInfo
4. update route: PUT /admin/route/:routeID          json body: RouteInfo
5. delete route: DELETE /admin/route/:routeID
*/

func init() {
	AddRegister(registerRoute)
}

type routeInfo struct {
	Id          string `json:"id"`
	ServiceName string `json:"serviceName"`
	MethodName  string `json:"methodName"`
	HttpMethod  string `json:"httpMethod"`
	Url         string `json:"url"`
}

func registerRoute(admin *route.RouterGroup) {
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
