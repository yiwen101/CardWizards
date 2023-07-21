package router

import (
	"github.com/cloudwego/hertz/pkg/app"
	//"github.com/yiwen101/CardWizards/service"
)

// router is responsible for finding the corresponding service and method according to the request path.
type RouteManager interface {
	ValidateRoute(c *app.RequestContext, httpMethod string) (string, string, error)
	InitRoute() error
	//UpdateRoute()
	GetRoute(path string) (string, string, bool)
	Get() map[string]api
}

/*
func GetRouteManager() (RouteManager, error) {
	if routeManager == nil {
		rm, err := newRouteManagerImpl()
		if err != nil {
			return nil, err
		}
		routeManager = rm
	}
	return routeManager, nil
} */

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

func (r *routeManagerImpl) InitRoute() error {
	if r.route != nil {
		return nil
	}

	r.route = make(map[string]api)
	services, err := r.dm.GetAllServiceNames()
	if err != nil {
		return err
	}
	for _, serviceName := range services {
		methods, err := r.dm.GetAllMethodNames(serviceName)
		if err != nil {
			return err
		}
		for _, methodName := range methods {
			url := "/" + serviceName + "/" + methodName
			api := api{methodName: methodName, serviceName: serviceName}
			r.route[url] = api
		}
	}
	return nil
}

func (r *routeManagerImpl) GetRoute(path string) (string, string, bool) {
	api, ok := r.route[path]
	if !ok {
		return "", "", false
	}
	return api.serviceName, api.methodName, true
}

func (r *routeManagerImpl) Get() map[string]api {
	return r.route
}
