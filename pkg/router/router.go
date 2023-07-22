package router

import (
	"fmt"

	"github.com/yiwen101/CardWizards/pkg/store"
	"github.com/yiwen101/CardWizards/pkg/utils"
)

// router is responsible for finding the corresponding service and method according to the request path.

func GetRoute(method, url string) (*RouteData, bool) {
	_, ok := localStore[method]
	if !ok {
		return nil, false
	}
	return localStore[method].Get(url)
}

func AddRoute(method, url string, data *RouteData) error {
	_, ok := localStore[method]
	if !ok {
		return fmt.Errorf("invalid http method, received %s", method)
	}
	return localStore[method].Add(url, data)
}

func DeleteRoute(method, url string) error {
	_, ok := localStore[method]
	if !ok {
		return fmt.Errorf("invalid http method, received %s", method)
	}
	return localStore[method].Delete(url)
}

func UpdateRoute(method, url, newMethod, newUrl string) error {
	data, ok := GetRoute(method, url)
	if !ok {
		return fmt.Errorf("route does not exist")
	}
	err := DeleteRoute(method, url)
	if err != nil {
		return err
	}
	return AddRoute(newMethod, newUrl, data)
}

var localStore map[string]*utils.MutexMap[string, *RouteData]

func init() {
	localStore = make(map[string]*utils.MutexMap[string, *RouteData])
	httpMethods := utils.HTTPMethods()
	for _, method := range httpMethods {
		localStore[method] = utils.NewMutexMap[string, *RouteData]()
	}

	store.InfoStore.RegisterApiRouteListener(&methodRouteHandeler{})
	store.InfoStore.RegisterServiceMapListener(&serviceRouteHandeler{})
}

type RouteData struct {
	MethodName  string
	ServiceName string
}

type methodRouteHandeler struct{}

func (rh *methodRouteHandeler) OnStatechanged(data ...interface{}) error {
	serviceName := data[0].(string)
	methodName := data[1].(string)
	url := data[2].(string)
	httpMethod := data[3].(string)
	isAdd := data[4].(bool)
	if isAdd {
		return AddRoute(httpMethod, url, &RouteData{methodName, serviceName})
	} else {
		return DeleteRoute(httpMethod, url)
	}
}

type serviceRouteHandeler struct{}

func (rh *serviceRouteHandeler) OnStatechanged(data ...interface{}) error {
	isAdd := data[0].(bool)
	meta := data[1].(*store.ServiceMeta)
	apis := meta.APIs
	for _, api := range apis {
		routeData := &RouteData{api.MethodName, api.ServiceName}
		for key, value := range api.Routes {
			for url := range value {
				if isAdd {
					err := AddRoute(key, url, routeData)
					if err != nil {
						return err
					}
				} else {
					err := DeleteRoute(key, url)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
