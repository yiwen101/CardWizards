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

/*
import (
	"bytes"
	"context"

	//"fmt"
	"log"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	//"github.com/cloudwego/hertz/pkg/app/client/loadbalance"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	loadbalance "github.com/cloudwego/kitex/pkg/loadbalance"

	"github.com/kitex-contrib/registry-nacos/resolver"
	handler "github.com/yiwen101/CardWizards/biz/handler"
)



func interface{} Serve(c *app.RequestContext) {
	serviceName := c.Param("serviceName")
	genericClient , err := GetGenericClient(serviceName)
	if err != nil {
		return &descriptor.HTTPResponse{
			StatusCode: http.StatusNotFound,
			Body:       []byte("service not found"),
		}
	}


}

				resolved, err := resolver.NewDefaultNacosResolver()
				if err != nil {
					panic(err)
				}
				// generic: first, make sure idl fulfill the requirement
				//read the idl file and generate the provider
				// todo: read all files in the IDL before running the server
				p, err := generic.NewThriftFileProvider("IDL/arithmatic.thrift")
				if err != nil {
					log.Println("error in generic.NewThriftFileProvider")
					panic(err)
				}
				//construct a generic httprequest with the provider
				g, err := generic.HTTPThriftGeneric(p)
				if err != nil {
					log.Println("error in generic.HTTPThriftGeneric")
					panic(err)
				}
				//make a generic client with the generic httprequest
				// but before that, add a load balance solution
				// todo: enable choice of load balance strategy? study how to add command line ->how: (service, strategy); is empty then all defaut, otherwise check before each time creating load balancer

				/*
				opt := loadbalance.NewConsistentHashOption(func(ctx context.Context, request interface{}) string {
					key, _ := ctx.Value(ctxConsistentKey).(string)
					return key
				})

				lb3 := loadbalance.NewWeightedRoundRobinBalancer()

				cli, err := genericclient.NewClient(
					serviceName,
					g,
					client.WithHostPorts("0.0.0.0:8889"),
					client.WithResolver(resolved),
					client.WithLoadBalancer(lb3),
				)

				// todo: general purpose handler

				if err != nil {
					log.Println("error in genericclient.NewClient")
					panic(err)
				}
				// optional: 构建一个http request
				//string(c.Request.URI().Path()) == /gateway/arith/add. get function lookup failed

				s := "/arith/add" + c.URI().QueryArgs().String()
				log.Println(s)
				httpReq, err := http.NewRequest(http.MethodGet, "/arith/add?"+c.URI().QueryArgs().String(), bytes.NewBuffer(c.Request.Body()))
				//log.Println("http request body is: ", c.Request.Body())
				//log.Println("http request url is: ", c.Request.RequestURI())

				if err != nil {
					log.Println("error in constructing http request")
					panic(err)
				}

				// 将http request转换成generic request
				customReq, err := generic.FromHTTPRequest(httpReq)
				if err != nil {
					log.Println("error in converting http request to generic request")
					panic(err)
				}

				// call
				genericResponse, err := cli.GenericCall(ctx, "", customReq)
				if err != nil {
					log.Println("error in generic call")
					panic(err)
				}

				resp := genericResponse.(*generic.HTTPResponse)
				// return response
				c.JSON(int(resp.StatusCode), resp.Body)

			} */
