package router

import (
	"context"
	"log"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yiwen101/CardWizards/common"
	"github.com/yiwen101/CardWizards/service"
)

type RouteManager interface {
	RegisterRoutes(r *server.Hertz)
}

type routeManagerImpl struct {
	routes []route
}

func NewRouteManager() RouteManager {
	return &routeManagerImpl{}
}

type route struct {
	httpMethod string
	route      string
	handler    func(ctx context.Context, c *app.RequestContext)
}

func (r *routeManagerImpl) generateRoutes() {
	r.routes = []route{}

	methods := []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead, http.MethodOptions}

	for _, method := range methods {
		handlerFunc := service.GenericHandlerFor(method)
		route := route{
			httpMethod: method,
			route:      common.RefaultRoute,
			handler:    handlerFunc,
		}
		r.routes = append(r.routes, route)
	}
	// todo
	// complex feature: use the annotation of the thrift file
	// intermediate feature: give route at command line
}

func (rm *routeManagerImpl) RegisterRoutes(r *server.Hertz) {
	rm.generateRoutes()
	for _, route := range rm.routes {
		switch route.httpMethod {
		case http.MethodGet:
			r.GET(route.route, route.handler)
			continue
		case http.MethodPost:
			r.POST(route.route, route.handler)
			continue
		case http.MethodPut:
			r.PUT(route.route, route.handler)
			continue
		case http.MethodDelete:
			r.DELETE(route.route, route.handler)
			continue
		case http.MethodPatch:
			r.PATCH(route.route, route.handler)
			continue
		case http.MethodHead:
			r.HEAD(route.route, route.handler)
			continue
		case http.MethodOptions:
			r.OPTIONS(route.route, route.handler)
			continue
		default:
			log.Println(" unsupported http method, invalid route ")
			continue
		}
	}
}
