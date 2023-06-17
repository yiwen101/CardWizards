package service

import (
	"bytes"
	"log"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
)

func HasService(serviceName string) bool {
	return serviceToClientMap[serviceName] != nil
}

func GetClient(serviceName string) (genericclient.Client, error) {

	return serviceToClientMap[serviceName], nil
}

func BuildRequest(c *app.RequestContext, method string) *descriptor.HTTPRequest {
	httpReq, err := http.NewRequest(method, c.Request.URI().String(), bytes.NewBuffer(c.Request.Body()))
	if err != nil {
		hlog.Info("error in constructing http request")
		panic(err)
	}

	// 将http request转换成generic request
	customReq, err := generic.FromHTTPRequest(httpReq)
	if err != nil {
		log.Println("error in converting http request to generic request")
		panic(err)
	}
	return customReq
}

