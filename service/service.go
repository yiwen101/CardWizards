package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yiwen101/CardWizards/router"
	client "github.com/yiwen101/CardWizards/service/clients"
	"github.com/yiwen101/CardWizards/service/validate"
)

func GenericHandlerFor(method string) func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {

		routeManager, err := router.GetRouteManager()
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error in getting the route manager: "+err.Error())
			return
		}

		serviceName, methodName, err := routeManager.ValidateRoute(c, method)
		if err != nil {
			c.String(http.StatusBadRequest, "invalid route"+err.Error())
		}

		validator := validate.NewValidator()
		err = validator.ValidateBody(ctx, c)

		if err != nil {
			c.String(http.StatusBadRequest, "invalid route: "+err.Error())
			return
		}

		cli, err := client.GetGenericClientforService(serviceName)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error in getting the client: "+err.Error())
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
	}
}

// is directly useful for httpGenericcall, but not for jsonGeneric call; nevertheless, is still needed for my validator to function
