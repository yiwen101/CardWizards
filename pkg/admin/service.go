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
List of APIs:
1. get all service names: GET /admin/service
2. get a service info: GET /admin/service/:serviceID
3. add service: POST /admin/service,                     json body: ServiceInfo
4. update service: PUT /admin/service/:serviceID,        json body: ServiceInfo
5. delete service: DELETE /admin/service/:serviceID
*/
func init() {
	AddRegister(registerService)
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

func registerService(admin *route.RouterGroup) {
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
