package service

import (
	"bytes"
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/yiwen101/CardWizards/temp"
	"github.com/yiwen101/CardWizards/validate"
)

func genericHandlerFor(method string) func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		serviceName := c.Param("serviceName")
		methodName := c.Param("methodName")

		validator := validate.NewValidator()

		err := validator.ValidateServiceMethodAndBody(ctx, c)

		if err != nil {
			return
		}
		cli, ok := temp.ServiceToClientMap[serviceName]
		if !ok {
			c.String(http.StatusInternalServerError, "Internal Server Error in getting the client")
			return
		}

		req, err := buildRequest(c, method)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}

		genericResponse, err := cli.GenericCall(ctx, methodName, req)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error in making the call")
			return
		}

		resp := genericResponse.(*generic.HTTPResponse)
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
