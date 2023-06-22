package router

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	desc "github.com/yiwen101/CardWizards/common/descriptor"
)

var routeManager RouteManager
var routes []Route

type Route struct {
	httpMethod  string
	ServiceName string
	MethodName  string
}

type routeManagerImpl struct {
	dm      desc.DescsManager
	cache   map[string]map[string]Route
	routers map[string]descriptor.Router
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
	if r.routers == nil {
		r.routers = r.dm.GetRouters()
	}

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
	if r.cache == nil {
		r.cache = make(map[string]map[string]Route)
	}
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

func newRouteManagerImpl() (RouteManager, error) {
	dm, err := desc.GetDescriptorManager()
	if err != nil {
		return nil, err
	}
	return &routeManagerImpl{dm: dm}, nil
}
