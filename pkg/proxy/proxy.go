package proxy

/*
 This package is responsible for providing the api gateway service
*/
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yiwen101/CardWizards/pkg/caller"
	"github.com/yiwen101/CardWizards/pkg/router"
	"github.com/yiwen101/CardWizards/pkg/store"
	"github.com/yiwen101/CardWizards/pkg/validator"
)

func Register(r *server.Hertz) {
	r.Any("proxy/*path",
		func(ctx context.Context, c *app.RequestContext) {
			Proxy.Serve(ctx, c, nil)
		})
}

// similar to a list of middleware
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

/*
to avoid repeating myself, provide a baseHandler struct to easily facilitate the creation of new handlerChain
with a func(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool) . If
The bool returned is false, the rest of the chain will not be executed.
Each handlerChain node should decide whether to abort the ctx with some code or message themselves.
*/
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
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Content-Type: "+str)
		return nil, false
	}
	return route, true
}

func routeHandler(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool) {
	path := "/" + string(c.Param("path"))
	method := string(c.Method())
	r, ok := router.GetRoute(method, path)
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, "404 page not found")
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
		c.AbortWithStatusJSON(http.StatusNotFound, "404 page not found")
		return nil, false
	}
	return route, meta.IsOn
}

func validationGateHandler(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool) {
	meta, err := store.InfoStore.CheckAPIStatus(route.ServiceName, route.MethodName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return nil, false
	}
	if meta.ValidationOn {
		good, errInfo := validator.Validate(c, route.ServiceName, route.MethodName)
		if !good {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid body: "+errInfo.Error())
			return nil, false
		}
	}

	return route, true
}

func mainHandler(ctx context.Context, c *app.RequestContext, route *router.RouteData) (*router.RouteData, bool) {
	client, ok := caller.GetClient(route.ServiceName)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "caller for the service not found")
		return nil, false
	}

	jsonbytes, err := c.Body()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Internal Server Error in marshalling the json body: "+err.Error())
		return nil, false
	}

	genericResponse, err := client.GenericCall(ctx, route.MethodName, string(jsonbytes))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, "GenericCall failed, error: "+err.Error())
		go log.Println(fmt.Sprintf("GenericCall failed for %s, error: %s", route.ServiceName, err.Error()))
		return nil, false
	}
	str, ok := genericResponse.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Internal Server Error in type assertion")
		return nil, false
	}

	// str is already Json string; no need to marshal
	c.String(200, str)
	c.SetContentType("application/json")

	return route, true
}
