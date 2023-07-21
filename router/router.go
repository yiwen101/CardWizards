package router

import (
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
)

// router is responsible for finding the corresponding service and method according to the request path.
type RouteManager interface {
	ValidateRoute(c *app.RequestContext, httpMethod string) (string, string, error)
	InitRoute() error
	AddRoute(method string, url string, newApi *Api) error
	UpdateRoute(method, url, newMethod, newUrl string) error
	DeleteRoute(method, path string) error
	GetRoute(method, path string) (*Api, bool)
	Get() map[string]*mutexMap
}

func (r *routeManagerImpl) ValidateRoute(c *app.RequestContext, httpMethod string) (string, string, error) {
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

	r.route = make(map[string]*mutexMap)

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
			api := Api{MethodName: methodName, ServiceName: serviceName, IsOn: true}
			err = r.AddRoute("GET", url, &api)
			if err != nil {
				log.Fatal(fmt.Errorf("when creating route for " + serviceName + " " + methodName + ": " + err.Error()))
			}
		}
	}
	return nil
}

func (r *routeManagerImpl) DeleteRoute(method, path string) error {
	mMap, ok := r.route[method]
	if !ok {
		return fmt.Errorf("route does not exist")
	}
	return mMap.delete(path)
}

func (r *routeManagerImpl) GetRoute(method, path string) (*Api, bool) {
	mMap, ok := r.route[method]
	if !ok {
		return nil, false
	}
	return mMap.get(path)
}

func (r *routeManagerImpl) UpdateRoute(method string, url string, newMethod string, newUrl string) error {
	api, ok := r.GetRoute(method, url)
	if !ok {
		return fmt.Errorf("route does not exist")
	}
	err := r.DeleteRoute(method, url)
	if err != nil {
		return err
	}
	return r.AddRoute(newMethod, newUrl, api)
}

func (r *routeManagerImpl) AddRoute(method string, url string, newApi *Api) error {
	mMap, ok := r.route[method]
	if !ok {
		mMap = newMutexMap()
		r.route[method] = mMap
	}
	return mMap.add(url, newApi)
}

func (r *routeManagerImpl) Get() map[string]*mutexMap {
	return r.route
}
