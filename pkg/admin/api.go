package admin

import (
	"context"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/yiwen101/CardWizards/pkg/store"
)

/*
list of admin api:
11. get all api info: GET /admin/api
12. get an api info: GET /admin/api/:apiID
13. get all api under a service: GET /admin/api/:serviceID
14. update api options: PUT /admin/api/:apiID                 json body: ApiInfo
*/

func init() {
	AddRegister(registerAPI)
}

type apiInfo struct {
	Id               string `json:"id"`
	ServiceName      string `json:"serviceName"`
	MethodName       string `json:"methodName"`
	IsSleeping       bool   `json:"isSleeping"`
	ValidationStatus bool   `json:"validationOn"`
}

func registerAPI(admin *route.RouterGroup) {
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
}

func translateAPI(api *store.ApiMeta) apiInfo {
	return apiInfo{api.ServiceName + "/" + api.MethodName, api.ServiceName, api.MethodName, !api.IsOn, api.ValidationOn}
}
