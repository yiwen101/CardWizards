package service

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	client "github.com/yiwen101/CardWizards/service/clients"
	"github.com/yiwen101/CardWizards/service/validate"
)

func GenericHandlerFor(method string) func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {

		serviceName := c.Param("serviceName")
		methodName := c.Param("methodName")

		req, err := buildRequest(c, method)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		validator := validate.NewValidator()

		serviceName, err = validator.ValidateRoute(serviceName, methodName, req)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		err = validator.ValidateBody(ctx, c)

		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		cli, ok := client.ServiceToClientMap[serviceName]
		if !ok {
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

		log.Println(genericResponse)
		resp, ok := genericResponse.(string)
		if !ok {
			c.String(http.StatusInternalServerError, "Internal Server Error in converting the generic response: "+err.Error())
			return
		}

		c.String(200, resp)
	}
}

// is directly useful for httpGenericcall, but not for jsonGeneric call; nevertheless, is still needed for my validator to function

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
