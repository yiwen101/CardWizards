package router

import (
	"github.com/cloudwego/hertz/pkg/app"
	//"github.com/yiwen101/CardWizards/service"
)

// router is responsible for finding the corresponding service and method according to the request path.
type RouteManager interface {
	ValidateRoute(c *app.RequestContext, httpMethod string) (string, string, error)
}

func GetRouteManager() (RouteManager, error) {
	if routeManager == nil {
		rm, err := newRouteManagerImpl()
		if err != nil {
			return nil, err
		}
		routeManager = rm
	}
	return routeManager, nil
}

func (r *routeManagerImpl) ValidateRoute(c *app.RequestContext, httpMethod string) (string, string, error) {
	serviceName, methodName := c.Param("serviceName"), c.Param("methodName")
	err := r.isGenericRoute(serviceName, methodName)
	if err == nil {
		return serviceName, methodName, nil
	}
	req, err := r.buildRequest(c, httpMethod)

	return r.isAnnotatedRoute(req)
}
