package router

import (
	"fmt"
	"sync"

	"github.com/yiwen101/CardWizards/pkg/store"
	"github.com/yiwen101/CardWizards/pkg/utils"
)

// router is responsible for finding the corresponding service and method according to the request path.

func GetRoute(method, url string) (*RouteData, bool) {
	_, ok := localStore[method]
	if !ok {
		return nil, false
	}
	return localStore[method].get(url)
}

func AddRoute(method, url string, data *RouteData) error {
	_, ok := localStore[method]
	if !ok {
		return fmt.Errorf("invalid http method, received %s", method)
	}
	return localStore[method].add(url, data)
}

func DeleteRoute(method, url string) error {
	_, ok := localStore[method]
	if !ok {
		return fmt.Errorf("invalid http method, received %s", method)
	}
	return localStore[method].delete(url)
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

var localStore map[string]*mutexMap

func init() {
	localStore = make(map[string]*mutexMap)
	httpMethods := utils.HTTPMethods()
	for _, method := range httpMethods {
		localStore[method] = newMutexMap()
	}

	store.InfoStore.RegisterApiRouteListener(&aPIrouteHandeler{})
	store.InfoStore.RegisterServiceMapListener(&serviceRouteHandeler{})
}

type mutexMap struct {
	// url -> arguements
	m   map[string]*RouteData
	mut sync.RWMutex
}

func newMutexMap() *mutexMap {
	return &mutexMap{m: make(map[string]*RouteData), mut: sync.RWMutex{}}
}

type RouteData struct {
	MethodName  string
	ServiceName string
}

func (m *mutexMap) get(key string) (*RouteData, bool) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	v, ok := m.m[key]
	return v, ok
}

func (m *mutexMap) delete(key string) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	_, ok := m.m[key]
	if !ok {
		return fmt.Errorf("route not found")
	}
	delete(m.m, key)
	return nil
}

func (m *mutexMap) add(key string, value *RouteData) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	_, ok := m.m[key]
	if ok {
		return fmt.Errorf("route already exists")
	}
	m.m[key] = value
	return nil
}

type aPIrouteHandeler struct{}

// serviceName, methodName, url, httpMethod, isAdd?
func (rh *aPIrouteHandeler) OnStatechanged(data ...interface{}) error {
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

// isAdd?, serviceMeta
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
