package proxyService

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yiwen101/CardWizards/pkg/caller"
	"github.com/yiwen101/CardWizards/pkg/router"
	"github.com/yiwen101/CardWizards/pkg/service/validator"
	"github.com/yiwen101/CardWizards/pkg/store"
)

var Proxy filter

func init() {
	porxyGateFilter := newFilter(proxyGateHandler)
	routeFilter := newFilter(routeHandler)
	apiGateFilter := newFilter(apiGateHandler)
	validationGateFilter := newFilter(validationGateHandler)
	hander := newFilter(mainHandler)

	Proxy = makeFilterChain(porxyGateFilter, routeFilter, apiGateFilter, validationGateFilter, hander)
}

type filter interface {
	Serve(ctx context.Context, c *app.RequestContext, route *router.RouteData)
	SetNext(f filter)
}

type baseFilter struct {
	hander func(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool)
	next   filter
}

func (f *baseFilter) SetNext(next filter) {
	f.next = next
}

func (f *baseFilter) Serve(ctx context.Context, c *app.RequestContext, route *router.RouteData) {
	if r, ok := f.hander(ctx, c, route); ok {
		if f.next != nil {
			f.next.Serve(ctx, c, r)
		}
	}
}

func newFilter(hander func(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool)) filter {
	return &baseFilter{hander: hander}
}

func proxyGateHandler(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool) {
	ok, err := store.InfoStore.CheckProxyStatus()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return nil, false
	}
	return nil, ok
}

func routeHandler(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool) {
	path := string(c.URI().Path())
	method := string(c.Method())
	return router.GetRoute(method, path)
}

func apiGateHandler(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool) {
	if route == nil {
		c.String(http.StatusInternalServerError, "router is ill")
		return nil, false
	}
	meta, err := store.InfoStore.CheckAPIStatus(route.ServiceName, route.MethodName)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return nil, false
	}
	return route, meta.IsOn
}

func validationGateHandler(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool) {
	meta, err := store.InfoStore.CheckAPIStatus(route.ServiceName, route.MethodName)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return nil, false
	}
	if meta.ValidationOn {
		good, err := validator.Validate(c, route.ServiceName, route.MethodName)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return nil, false
		}
		if !good {
			c.String(http.StatusBadRequest, "Invalid body: "+err.Error())
			return nil, false
		}
	}

	return route, true
}

// posible extension: enabling managing call options on api level; should add to filter
func mainHandler(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool) {
	client, ok := caller.GetClient(route.ServiceName)
	if !ok {
		c.String(http.StatusInternalServerError, "caller for the service not found")
		return nil, false
	}

	jsonbytes, err := c.Body()
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error in marshalling the json body: "+err.Error())
		return nil, false
	}

	genericResponse, err := client.GenericCall(ctx, route.MethodName, string(jsonbytes))
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error in making the call: "+err.Error())
		return nil, false
	}

	resp, ok := genericResponse.(string)
	if !ok {
		c.String(http.StatusInternalServerError, "Internal Server Error in converting the generic response: "+err.Error())
		return nil, false
	}

	c.JSON(200, resp)

	return route, true

}

func makeFilterChain(f ...filter) filter {
	if len(f) == 0 {
		return nil
	}
	for i := 0; i < len(f)-1; i++ {
		f[i].SetNext(f[i+1])
	}
	return f[0]
}
