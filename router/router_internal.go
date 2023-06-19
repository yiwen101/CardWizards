package router

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	desc "github.com/yiwen101/CardWizards/common/descriptor"
)

func (r *routeManagerImpl) isGenericRoute(serviceName, methodName string) error {

	_, err := r.dm.GetServiceDescriptor(serviceName)
	if err != nil {
		return fmt.Errorf("service %s not found", serviceName)
	}
	_, err = r.dm.GetFunctionDescriptor(serviceName, methodName)
	return err
}

func (r *routeManagerImpl) isAnnotatedRoute(req *descriptor.HTTPRequest) (string, string, error) {
	return r.dm.GetMathchedRouterName(req)
}

func (r *routeManagerImpl) buildRequest(c *app.RequestContext, method string) (*descriptor.HTTPRequest, error) {
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

func newRouteManagerImpl() (RouteManager, error) {
	dm, err := desc.GetDescriptorManager()
	if err != nil {
		return nil, err
	}
	return &routeManagerImpl{dm: dm}, nil
}

var routeManager RouteManager

type routeManagerImpl struct {
	dm desc.DescsManager
}
