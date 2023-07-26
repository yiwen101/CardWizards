package proxy

// service is responsible for providing the api gateway service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yiwen101/CardWizards/pkg/caller"
	"github.com/yiwen101/CardWizards/pkg/proxy/validator"
	"github.com/yiwen101/CardWizards/pkg/router"
	"github.com/yiwen101/CardWizards/pkg/store"
)

func Register(r *server.Hertz) {

	r.Any("/*path", func(ctx context.Context, c *app.RequestContext) { Proxy.Serve(ctx, c, nil) })
}

var Proxy handlerChain

func init() {
	porxyGate := newHandlerChainNode(proxyGateHandler)
	checkContentType := newHandlerChainNode(contentTypeHandler)
	router := newHandlerChainNode(routeHandler)
	apiGate := newHandlerChainNode(apiGateHandler)
	validator := newHandlerChainNode(validationGateHandler)
	mainHander := newHandlerChainNode(mainHandler)

	Proxy = makeFilterChain(porxyGate, checkContentType, router, apiGate, validator, mainHander)
}

type handlerChain interface {
	Serve(ctx context.Context, c *app.RequestContext, route *router.RouteData)
	SetNext(f handlerChain)
}

type baseHandler struct {
	hander func(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool)
	next   handlerChain
}

func (f *baseHandler) SetNext(next handlerChain) {
	f.next = next
}

func (f *baseHandler) Serve(ctx context.Context, c *app.RequestContext, route *router.RouteData) {
	if r, ok := f.hander(ctx, c, route); ok {
		if f.next != nil {
			f.next.Serve(ctx, c, r)
		}
	}
}

func newHandlerChainNode(hander func(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool)) handlerChain {
	return &baseHandler{hander: hander}
}

func makeFilterChain(f ...handlerChain) handlerChain {
	if len(f) == 0 {
		return nil
	}
	for i := 0; i < len(f)-1; i++ {
		f[i].SetNext(f[i+1])
	}
	return f[0]
}

func proxyGateHandler(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool) {
	ok, err := store.InfoStore.CheckProxyStatus()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return nil, false
	}
	if !ok {
		c.Abort()
	}
	return nil, ok
}

func contentTypeHandler(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool) {
	bytes := c.ContentType()
	str := string(bytes)

	if !strings.Contains(str, "application/json") {
		c.String(http.StatusBadRequest, "Invalid Content-Type: "+str)
		return nil, false
	}
	return route, true
}

func routeHandler(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool) {
	path := string(c.URI().Path())
	method := string(c.Method())
	r, ok := router.GetRoute(method, path)
	if !ok {
		c.String(http.StatusNotFound, "404 page not found")
		return nil, false
	}
	return r, true
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
	if !meta.IsOn {
		c.String(http.StatusNotFound, "404 page not found")
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
		good, errInfo := validator.Validate(c, route.ServiceName, route.MethodName)
		if !good {
			c.String(http.StatusBadRequest, "Invalid body: "+errInfo.Error())
			return nil, false
		}
	}

	return route, true
}

// posible extension: enabling managing call options per api level; should add to filter
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
		c.String(http.StatusBadGateway, "GenericCall failed, error: "+err.Error())
		log.Println(fmt.Sprintf("GenericCall failed for %s, error: %s", route.ServiceName, err.Error()))
		return nil, false
	}
	str, ok := genericResponse.(string)
	if !ok {
		c.String(http.StatusInternalServerError, "Internal Server Error in type assertion")
		return nil, false
	}
	c.String(200, str)
	c.SetContentType("application/json")

	return route, true
}
