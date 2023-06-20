package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/yiwen101/CardWizards/router"
	client "github.com/yiwen101/CardWizards/service/clients"
	"github.com/yiwen101/CardWizards/service/validate"
)

type HandlerManager interface {
	HandlerForAnnotatedRoutes(httpMethod string) (func(ctx context.Context, c *app.RequestContext), error)
	HandlerForRoute(serviceName, methodName string) (func(ctx context.Context, c *app.RequestContext), error)
}

func GetHandlerManager() (HandlerManager, error) {
	if hm == nil {
		hm = &handlerManagerImpl{}
	}
	return hm, nil
}
func (hm *handlerManagerImpl) HandlerForAnnotatedRoutes(httpMethod string) (func(ctx context.Context, c *app.RequestContext), error) {
	routeManager, err := router.GetRouteManager()
	if err != nil {
		hlog.Fatal("Internal Server Error in getting the route manager: ", err)
		return nil, err
	}

	handlerCache := &handlerCache{}

	return func(ctx context.Context, c *app.RequestContext) {

		serviceName, methodName, err := routeManager.ValidateRoute(c, httpMethod)
		if err != nil {
			c.String(http.StatusBadRequest, "invalid route: "+err.Error())
			return
		}

		handler, ok := handlerCache.get(serviceName, methodName)
		if ok {
			handler(ctx, c)
			return
		}

		handler, err = hm.HandlerForRoute(serviceName, methodName)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error in getting the handler: "+err.Error())
			return
		}
		handlerCache.save(serviceName, methodName, handler)
		handler(ctx, c)
	}, nil
}

func (hm *handlerManagerImpl) HandlerForRoute(serviceName, methodName string) (func(ctx context.Context, c *app.RequestContext), error) {

	cli, err := client.GetGenericClientforService(serviceName)
	if err != nil {
		hlog.Fatal("Internal Server Error in getting the client: ", err)
		return nil, err
	}

	validator, err := validate.NewValidatorFor(serviceName, methodName)
	if err != nil {
		hlog.Fatal("Internal Server Error in getting the validator: ", err)
	}

	return func(ctx context.Context, c *app.RequestContext) {

		err = validator.ValidateBody(c, serviceName, methodName)

		if err != nil {
			c.String(http.StatusBadRequest, "invalid route: "+err.Error())
			return
		}

		jsonbytes, err := c.Body()
		// should not happen, otherwise indicate that there is problem with my validator
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error in marshalling the json body: "+err.Error())
			return
		}
		var jsonString string
		json.Unmarshal(jsonbytes, &jsonString)

		c.String(http.StatusOK, jsonString)

		genericResponse, err := cli.GenericCall(ctx, methodName, jsonString)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error in making the call: "+err.Error())
			return
		}

		resp, ok := genericResponse.(string)
		if !ok {
			c.String(http.StatusInternalServerError, "Internal Server Error in converting the generic response: "+err.Error())
			return
		}

		c.JSON(200, resp)
	}, nil
}

type handlerCache struct {
	m map[string]map[string]func(ctx context.Context, c *app.RequestContext)
}

func (hc *handlerCache) get(serviceName, methodName string) (func(ctx context.Context, c *app.RequestContext), bool) {
	if hc.m == nil {
		hc.m = make(map[string]map[string]func(ctx context.Context, c *app.RequestContext))
	}
	if hc.m[serviceName] == nil {
		hc.m[serviceName] = make(map[string]func(ctx context.Context, c *app.RequestContext))
	}
	handler, ok := hc.m[serviceName][methodName]
	return handler, ok
}

func (hc *handlerCache) save(serviceName, methodName string, handler func(ctx context.Context, c *app.RequestContext)) {
	if hc.m == nil {
		hc.m = make(map[string]map[string]func(ctx context.Context, c *app.RequestContext))
	}
	if hc.m[serviceName] == nil {
		hc.m[serviceName] = make(map[string]func(ctx context.Context, c *app.RequestContext))
	}
	hc.m[serviceName][methodName] = handler
}

type handlerManagerImpl struct{}

var hm HandlerManager

// is directly useful for httpGenericcall, but not for jsonGeneric call; nevertheless, is still needed for my validator to function
