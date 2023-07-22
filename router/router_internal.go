package router

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/yiwen101/CardWizards/pkg/store"
	desc "github.com/yiwen101/CardWizards/pkg/store/descriptor"
)

var routeManager *routeManagerImpl

type Route struct {
	httpMethod  string
	ServiceName string
	MethodName  string
}

type Api struct {
	MethodName  string
	ServiceName string
	IsOn        bool
}

type routeManagerImpl struct {
	store *store.Store
	dm    desc.DescsManager
	// to update?
	cache   map[string]map[string]Route
	routers map[string]descriptor.Router
	route   map[string]*mutexMap
}

type mutexMap struct {
	m   map[string]*Api
	mut sync.RWMutex
}

func (m *mutexMap) get(key string) (*Api, bool) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	v, ok := m.m[key]
	return v, ok
}

func (m *mutexMap) delete(key string) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	_, ok := m.m[key]
	if !ok {
		return fmt.Errorf("route not found")
	}
	delete(m.m, key)
	return nil
}

func (m *mutexMap) add(key string, value *Api) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	_, ok := m.m[key]
	if ok {
		return fmt.Errorf("route already exists")
	}
	m.m[key] = value
	return nil
}

func newMutexMap() *mutexMap {
	return &mutexMap{m: make(map[string]*Api), mut: sync.RWMutex{}}
}

func (d *Route) GetRoute() (httpMethod string, path string) {
	path = "/" + d.ServiceName + "/" + d.MethodName
	return d.httpMethod, path
}

/*
func (r *routeManagerImpl) isGenericRoute(serviceName, methodName string) error {

	_, err := r.dm.GetServiceDescriptor(serviceName)
	if err != nil {
		return fmt.Errorf("service %s not found", serviceName)
	}
	_, err = r.dm.GetFunctionDescriptor(serviceName, methodName)
	return err
}
*/

func (r *routeManagerImpl) isAnnotatedRoute(req *descriptor.HTTPRequest) (string, string, error) {
	serviceName, methodName, ok := r.findInCache(req)
	if ok {
		return serviceName, methodName, nil
	}

	for serviceName, router := range r.routers {
		des, err := router.Lookup(req)
		if err == nil {
			r.saveInCache(req.Method, req.Path, serviceName, des.Name)
			return serviceName, des.Name, nil
		}
	}

	return "", "", fmt.Errorf("service not found")
}

func (r *routeManagerImpl) findInCache(req *descriptor.HTTPRequest) (string, string, bool) {
	httpMehtod, path := req.Method, req.Path
	m, ok := r.cache[httpMehtod]
	if ok {
		if pair, ok := m[path]; ok {
			return pair.ServiceName, pair.MethodName, true
		}
	}
	return "", "", false
}

func (r *routeManagerImpl) saveInCache(httpMehtod, path, serviceName, methodName string) {
	m, ok := r.cache[httpMehtod]
	if !ok {
		m = make(map[string]Route)
		r.cache[httpMehtod] = m
	}
	m[path] = Route{ServiceName: serviceName, MethodName: methodName}
}

func (r *routeManagerImpl) buildRequest(c *app.RequestContext, method string) (*descriptor.HTTPRequest, error) {
	httpReq, err := http.NewRequest(method, c.Request.URI().String(), bytes.NewBuffer(c.Request.Body()))
	if err != nil {
		return nil, err
	}

	// 将http request转换成generic request
	customReq, err := generic.FromHTTPRequest(httpReq)
	if err != nil {
		return nil, err
	}
	return customReq, nil
}

func GetRouteManager() (RouteManager, error) {
	if routeManager != nil {
		return routeManager, nil
	}

	dm, err := desc.GetDescriptorManager()
	if err != nil {
		return nil, err
	}
	cache := make(map[string]map[string]Route)
	routers := dm.GetRouters()
	routeManager = &routeManagerImpl{dm: dm, store: nil, cache: cache, routers: routers, route: nil}
	return routeManager, nil
}
