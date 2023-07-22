package router

/*
import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	desc "github.com/yiwen101/CardWizards/pkg/store/descriptor"
)

var cache map[string]map[string]*RouteData
var routers map[string]descriptor.Router

func init() {
	cache = make(map[string]map[string]*RouteData)
	routers = make(map[string]descriptor.Router)
}

func ValidateRoute(c *app.RequestContext, httpMethod string) (*RouteData, bool) {
	req, err := buildRequest(c, httpMethod)
	if err != nil {
		return nil, false
	}
	return isAnnotatedRoute(req)
}

func isAnnotatedRoute(req *descriptor.HTTPRequest) (string, string, error) {
	serviceName, methodName, ok := findInCache(req)
	if ok {
		return serviceName, methodName, nil
	}

	for serviceName, router := range routers {
		des, err := router.Lookup(req)
		if err == nil {
			r.saveInCache(req.Method, req.Path, serviceName, des.Name)
			return serviceName, des.Name, nil
		}
	}

	return "", "", fmt.Errorf("service not found")
}

func findInCache(req *descriptor.HTTPRequest) (string, string, bool) {
	httpMehtod, path := req.Method, req.Path
	m, ok := cache[httpMehtod]
	if ok {
		if pair, ok := m[path]; ok {
			return pair.ServiceName, pair.MethodName, true
		}
	}
	return "", "", false
}

func saveInCache(httpMehtod, path, serviceName, methodName string) {
	m, ok := cache[httpMehtod]
	if !ok {
		m = make(map[string]Route)
		r.cache[httpMehtod] = m
	}
	m[path] = Route{ServiceName: serviceName, MethodName: methodName}
}

func buildRequest(c *app.RequestContext, method string) (*descriptor.HTTPRequest, error) {
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
*/
