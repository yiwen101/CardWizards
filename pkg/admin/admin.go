package admin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yiwen101/CardWizards/pkg/service"
	"github.com/yiwen101/CardWizards/pkg/store"
)

func RegisterProxy(r *server.Hertz) {

	r.Any("/*path", func(ctx context.Context, c *app.RequestContext) { service.Proxy.Serve(ctx, c, nil) })
}

func RegisterAdmin(r *server.Hertz) {
	admin := r.Group("/admin")
	admin.Use(func(ctx context.Context, c *app.RequestContext) {
		b, _ := c.Body()
		var j map[string]interface{}

		err := json.Unmarshal(b, &j)
		if err != nil {
			return false, err
		}
	})
	admin.GET("/allServices",
		func(ctx context.Context, c *app.RequestContext) {

			m, err := store.InfoStore.GetAllServiceNames()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, m)
		})
	admin.GET("/serviceInfo/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.GetServiceInfo(c.Param("serviceName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, s)
		})
	admin.GET("/apiInfo/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.CheckAPIStatus(c.Param("serviceName"), c.Param("methodName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, s)
		})
	admin.GET("/allAPIs/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.GetAPIs(c.Param("serviceName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, s)
		})
	admin.GET("/routes/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.GetRoutes(c.Param("serviceName"), c.Param("methodName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, s)
		})
	admin.GET("/proxyStatus",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.CheckProxyStatus()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, s)
		})
	admin.GET("/lbType/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			s, err := store.InfoStore.GetLbType(c.Param("serviceName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, s)
		})

	admin.POST("/turnOnProxy",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.TurnOnProxy()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Proxy is on")
		})
	admin.POST("/turnOffProxy",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.TurnOffProxy()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Proxy is off")
		})
	admin.POST("/turnOnService/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.TurnOnService(c.Param("serviceName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Service is on")
		})
	admin.POST("/turnOffService/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.TurnOffService(c.Param("serviceName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Service is off")
		})
	admin.POST("/turnOnAPI/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.TurnOnAPI(c.Param("serviceName"), c.Param("methodName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "API is on")
		})
	admin.POST("/turnOffAPI/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.TurnOffAPI(c.Param("serviceName"), c.Param("methodName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "API is off")
		})
	admin.POST("/turnOnValidation/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.TurnOnValidation(c.Param("serviceName"), c.Param("methodName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Validation is on")
		})
	admin.POST("/turnOffValidation/:serviceName/:methodName",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.TurnOffValidation(c.Param("serviceName"), c.Param("methodName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Validation is off")
		})

	admin.POST("/addService/:idlFileName/:clusterName",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.AddService(c.Param("idlFileName"), c.Param("clusterName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Service is added")
		})
	admin.POST("/updateService/:serviceName/:idlFileName/:clusterName",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.UpdateService(c.Param("serviceName"), c.Param("idlFileName"), c.Param("clusterName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Service is updated")
		})
	admin.POST("/removeService/:serviceName",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.RemoveService(c.Param("serviceName"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Service is removed")
		})
	admin.POST("/addRoute/:serviceName/:methodName/:url/:httpMethod",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.AddRoute(c.Param("serviceName"), c.Param("methodName"), c.Param("url"), c.Param("httpMethod"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Route is added")
		})
	admin.POST("/updateRoute/:serviceName/:methodName/:url/:httpMethod/:newUrl/:newHttpMethod",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.UpdateRoute(c.Param("serviceName"), c.Param("methodName"), c.Param("url"), c.Param("httpMethod"), c.Param("newUrl"), c.Param("newHttpMethod"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Route is updated")
		})
	admin.POST("/removeRoute/:serviceName/:methodName/:url/:httpMethod",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.RemoveRoute(c.Param("serviceName"), c.Param("methodName"), c.Param("url"), c.Param("httpMethod"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "Route is removed")
		})
	admin.POST("/setLbType/:serviceName/:lbType",
		func(ctx context.Context, c *app.RequestContext) {
			err := store.InfoStore.SetLbType(c.Param("serviceName"), c.Param("lbType"))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.String(http.StatusOK, "LbType is set")
		})
}
