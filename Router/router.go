package router

import (
	"github.com/cloudwego/hertz/pkg/app"
	//"github.com/yiwen101/CardWizards/service"
)

// router is responsible for finding the corresponding service and method according to the request path.
type RouteManager interface {
	ValidateRoute(c *app.RequestContext, httpMethod string) (string, string, error)
	GetRoutes() ([]Route, error)
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
	//serviceName, methodName := c.Param("serviceName"), c.Param("methodName")
	/*
		err := r.isGenericRoute(serviceName, methodName)
		if err == nil {
			return serviceName, methodName, nil
		}
	*/

	req, err := r.buildRequest(c, httpMethod)
	if err != nil {
		return "", "", err
	}
	return r.isAnnotatedRoute(req)
}

func (r *routeManagerImpl) GetRoutes() ([]Route, error) {
	if routes != nil {
		return routes, nil
	}
	services, err := r.dm.GetAllServiceNames()
	if err != nil {
		return nil, err
	}
	for _, serviceName := range services {
		methods, err := r.dm.GetAllMethodNames(serviceName)
		if err != nil {
			return nil, err
		}
		for _, methodName := range methods {
			route := Route{
				httpMethod:  "POST",
				ServiceName: serviceName,
				MethodName:  methodName,
			}
			routes = append(routes, route)
		}
	}
	return routes, nil
}
