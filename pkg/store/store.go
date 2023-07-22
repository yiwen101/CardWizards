package store

import (
	"sync"

	"github.com/cloudwego/kitex/pkg/generic/descriptor"
)

type proxyAdmin interface {
	GetAllServiceNames() (map[string]*ServiceMeta, error)
	CheckProxyStatus() (bool, error)
	TurnOnProxy() error  //proxyGate
	TurnOffProxy() error //proxyGate
	GetProxyAddress() (string, error)
	GetStoreAddress() (string, error)

	GetAPIs(serviceName string) (map[string]*ApiMeta, error)
	AddService(idlFileName, clusterName string) error                 //validator, client, descriptor, router
	UpdateService(serviceName, idlFileName, clusterName string) error //validator, client, descriptor, router
	RemoveService(serviceName string) error                           //validator, client, descriptor, router
	TurnOnService(serviceName string) error
	TurnOffService(serviceName string) error
	GetServiceInfo(serviceName string) (*ServiceMeta, error)

	CheckAPIStatus(serviceName, methodName string) (*ApiMeta, error)
	TurnOnAPI(serviceName, methodName string) error                 //APIGate
	TurnOffAPI(serviceName, methodName string) error                //APIGate
	TurnOnValidation(serviceName, methodName string) error          //validator
	TurnOffValidation(serviceName, methodName string) error         //validator
	AddRoute(serviceName, methodName, url, httpMethod string) error //router
	ModifyRoute(url, httpMethod, newUrl, newMethod string) error    //router
	RemoveRoute(url, httpMethod string) error                       //router
	GetRoutes(serviceName, methodName string) (map[string]map[string]bool, error)
	GetLbType(serviceName string) (string, error) //lb
	SetLbType(serviceName, lbType string) error   //lb
}

var InfoStore *Store

func load() {
	InfoStore = &Store{
		IsOn:                true,
		proxyStateListeners: []EventListener{},
		ServicesMap:         map[string]*ServiceMeta{},
		ServiceMapListeners: []EventListener{},
		ProxyAddress:        ""
		StoreAddress:        ""
		IdlFolderRelativePath: ""
	}


type Store struct {
	mutex sync.RWMutex

	IsOn                bool
	proxyStateListeners []EventListener
	ServicesMap         map[string]*ServiceMeta
	ServiceMapListeners []EventListener

	ProxyAddress          string
	StoreAddress          string
	IdlFolderRelativePath string
}

type ServiceMeta struct {
	ServiceName string
	ClusterName string

	Descriptor descriptor.ServiceDescriptor
	LbType     string

	lbStateListners []EventListener

	Apis map[string]*ApiMeta
}

type ApiMeta struct {
	ServiceName string
	MethodName  string

	apiStateListeners      []EventListener
	apiValidationListeners []EventListener
	apiRouteListeners      []EventListener

	ValidationOn bool
	// map[httpmethod]map[url]
	Routes map[string]map[string]bool
	IsOn   bool
}

func (s *Store) RegisterProxyStateListener(listener EventListener) {
	s.proxyStateListeners = append(s.proxyStateListeners, listener)
}
func (s *Store) RegisterServiceMapListener(listener EventListener) {
	s.ServiceMapListeners = append(s.ServiceMapListeners, listener)
}
func (s *Store) RegisterApiStateListener(serviceName, methodName string, listener EventListener) {
	s.ServicesMap[serviceName].Apis[methodName].apiStateListeners = append(s.ServicesMap[serviceName].Apis[methodName].apiStateListeners, listener)
}
func (s *Store) RegisterApiValidationListener(serviceName, methodName string, listener EventListener) {
	s.ServicesMap[serviceName].Apis[methodName].apiValidationListeners = append(s.ServicesMap[serviceName].Apis[methodName].apiValidationListeners, listener)
}
func (s *Store) RegisterApiRouteListener(serviceName, methodName string, listener EventListener) {
	s.ServicesMap[serviceName].Apis[methodName].apiRouteListeners = append(s.ServicesMap[serviceName].Apis[methodName].apiRouteListeners, listener)
}

type EventListener interface {
	OnStatechanged(data ...interface{})
}

func notifyStatechange(listeners []EventListener, data ...interface{}) {
	for _, listener := range listeners {
		listener.OnStatechanged(data...)
	}
}

func (s *Store) GetAllServiceNames() (map[string]*ServiceMeta, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.ServicesMap, nil
}

func (s *Store) CheckProxyStatus() (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.IsOn, nil
}

func (s *Store) TurnOnProxy() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.IsOn = true
	notifyStatechange(s.proxyStateListeners, true)
	return nil
}

func (s *Store) TurnOffProxy() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.IsOn = false
	notifyStatechange(s.proxyStateListeners, false)
	return nil
}

func (s *Store) GetProxyAddress() (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.ProxyAddress, nil
}

func (s *Store) GetStoreAddress() (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.StoreAddress, nil
}

func (s *Store) GetAPIs(serviceName string) (map[string]*ApiMeta, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.ServicesMap[serviceName].Apis, nil
}

func (s *Store) AddService(idlFileName, clusterName string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	dk, err := buildDescriptorKeeperFromPath(idlFileName, s.IdlFolderRelativePath)
	if err != nil {
		return err
	}
	sd, err := dk.get()
	if err != nil {
		return err
	}
	serviceName := sd.Name

	result := map[string]*ApiMeta{}

	for methodName := range sd.Functions {
		route := make(map[string]map[string]bool)
		route["GET"] = make(map[string]bool)
		route["GET"]["/"+serviceName+"/"+methodName] = true
		api := ApiMeta{
			ServiceName:  serviceName,
			MethodName:   methodName,
			ValidationOn: true,
			Routes:       route,
			IsOn:         true,
		}
		result[methodName] = &api
	}

	s.ServicesMap[serviceName] = &ServiceMeta{
		ServiceName: serviceName,
		ClusterName: clusterName,
		Descriptor:  *sd,
		//todo
		LbType: "random",
		Apis:   result,
	}

	notifyStatechange(s.ServiceMapListeners, serviceName)
	return nil
}

func (s *Store) RemoveService(serviceName string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	delete(s.ServicesMap, serviceName)
	notifyStatechange(s.ServiceMapListeners, serviceName)
	return nil
}

func (s *Store) UpdateService(serviceName, idlFileName, clusterName string) error {
	s.RemoveService(serviceName)
	s.AddService(idlFileName, clusterName)
	return nil
}

func (s *Store) TurnOnService(serviceName string) error {
	for methodName := range s.ServicesMap[serviceName].Apis {
		s.TurnOnAPI(serviceName, methodName)
	}
	return nil
}

func (s *Store) TurnOffService(serviceName string) error {
	for methodName := range s.ServicesMap[serviceName].Apis {
		s.TurnOffAPI(serviceName, methodName)
	}
	return nil
}

func (s *Store) GetServiceInfo(serviceName string) (*ServiceMeta, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.ServicesMap[serviceName], nil
}

func (s *Store) CheckAPIStatus(serviceName, methodName string) (*ApiMeta, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.ServicesMap[serviceName].Apis[methodName], nil
}

func (s *Store) TurnOnAPI(serviceName, methodName string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ServicesMap[serviceName].Apis[methodName].IsOn = true
	notifyStatechange(s.ServicesMap[serviceName].Apis[methodName].apiStateListeners, serviceName, methodName, true)
	return nil
}

func (s *Store) TurnOffAPI(serviceName, methodName string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ServicesMap[serviceName].Apis[methodName].IsOn = false
	notifyStatechange(s.ServicesMap[serviceName].Apis[methodName].apiStateListeners, serviceName, methodName, false)
	return nil
}

func (s *Store) TurnOnValidation(serviceName, methodName string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ServicesMap[serviceName].Apis[methodName].ValidationOn = true
	notifyStatechange(s.ServicesMap[serviceName].Apis[methodName].apiValidationListeners, serviceName, methodName, true)
	return nil
}

func (s *Store) TurnOffValidation(serviceName, methodName string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ServicesMap[serviceName].Apis[methodName].ValidationOn = false
	notifyStatechange(s.ServicesMap[serviceName].Apis[methodName].apiValidationListeners, serviceName, methodName, false)
	return nil
}

func (s *Store) AddRoute(serviceName, methodName, url, httpMethod string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ServicesMap[serviceName].Apis[methodName].Routes[httpMethod][url] = true
	notifyStatechange(s.ServicesMap[serviceName].Apis[methodName].apiRouteListeners, serviceName, methodName, url, httpMethod, true)
	return nil
}

func (s *Store) DeleteRoute(serviceName, methodName, url, httpMethod, newUrl, newMethod string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.ServicesMap[serviceName].Apis[methodName].Routes[httpMethod], url)
	notifyStatechange(s.ServicesMap[serviceName].Apis[methodName].apiRouteListeners, serviceName, methodName, url, httpMethod, false)
	return nil
}

func (s *Store) ModifyRoute(serviceName, methodName, url, httpMethod, newUrl, newMethod string) error {
	err := s.DeleteRoute(serviceName, methodName, url, httpMethod, newUrl, newMethod)
	if err != nil {
		return err
	}
	return s.AddRoute(serviceName, methodName, newUrl, newMethod)
}

func (s *Store) GetRoutes(serviceName, methodName string) (map[string]map[string]bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.ServicesMap[serviceName].Apis[methodName].Routes, nil
}

func (s *Store) GetLbType(serviceName string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.ServicesMap[serviceName].LbType, nil
}

func (s *Store) SetLbType(serviceName, lbType string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ServicesMap[serviceName].LbType = lbType
	notifyStatechange(s.ServicesMap[serviceName].lbStateListners, serviceName, lbType)
	return nil
}
