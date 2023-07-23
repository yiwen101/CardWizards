package admin

import (
	"context"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
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