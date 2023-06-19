package service

import (
	"bytes"
	"context"
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

		genericResponse, err := cli.GenericCall(ctx, methodName, req)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error in making the call: "+err.Error())
			return
		}

		resp, ok := genericResponse.(*generic.HTTPResponse)
		if !ok {
			c.String(http.StatusInternalServerError, "Internal Server Error in converting the generic response: "+err.Error())
			return
		}

		c.JSON(int(resp.StatusCode), resp.Body)
	}
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
