package store

import (
	"fmt"
	"log"
	"os"
	"sync"
)

//todo improve performance by adding multiple mutexes

type Admin interface {
	GetAllServiceNames() (map[string]*ServiceMeta, error)
	CheckProxyStatus() (bool, error)
	TurnOnProxy() error  //proxyGate
	TurnOffProxy() error //proxyGate

	GetAPIs(serviceName string) (map[string]*ApiMeta, error)
	AddService(idlFileName, clusterName string) error                 //validator, caller,  router, APIGate
	UpdateService(serviceName, idlFileName, clusterName string) error //validator, caller,  router, APIGate
	RemoveService(serviceName string) error                           //validator, caller,  router, APIGate
	TurnOnService(serviceName string) error
	TurnOffService(serviceName string) error
	GetServiceInfo(serviceName string) (*ServiceMeta, error)

	CheckAPIStatus(serviceName, methodName string) (*ApiMeta, error)
	TurnOnAPI(serviceName, methodName string) error                                       //APIGate
	TurnOffAPI(serviceName, methodName string) error                                      //APIGate
	TurnOnValidation(serviceName, methodName string) error                                //validator
	TurnOffValidation(serviceName, methodName string) error                               //validator
	AddRoute(serviceName, methodName, httpMethod, url string) error                       //router
	ModifyRoute(serviceName, methodName, httpMethod, url, newMethod, newUrl string) error //router
	RemoveRoute(serviceName, methodName, httpMethod, url string) error                    //router
	GetRoutes(serviceName, methodName string) (map[string]map[string]bool, error)
	GetLbType(serviceName string) (string, error) //caller
	SetLbType(serviceName, lbType string) error   //caller
}

var InfoStore *Store

func init() {

	InfoStore = &Store{
		mutex:                  sync.RWMutex{},
		IsOn:                   true,
		proxyStateListeners:    []EventListener{},
		ServiceMapListeners:    []EventListener{},
		apiStateListeners:      []EventListener{},
		apiValidationListeners: []EventListener{},
		apiRouteListeners:      []EventListener{},
		lbStateListners:        []EventListener{},
		ServicesMap:            map[string]*ServiceMeta{},
		ProxyAddress:           "",
		IdlFolderRelativePath:  "../../IDL",
	}
}

// default cluster name is the name of the idl file
func (s *Store) Load(ProxyAddress, IdlFolderRelativePath, password string) {
	s.ProxyAddress = ProxyAddress
	s.IdlFolderRelativePath = IdlFolderRelativePath
	s.Password = password
	thiriftFiles, err := os.ReadDir(s.IdlFolderRelativePath)
	if err != nil {
		log.Fatal(err)
	}

	waitGroup := sync.WaitGroup{}

	for _, file := range thiriftFiles {
		log.Printf("reading file : %s", file.Name())
		if file.IsDir() {
			log.Fatal("failure reading thrrift files at IDL directory as it contains directory")
		}

		if file.Name()[len(file.Name())-7:] != ".thrift" {
			log.Fatal("failure reading thrrift files at IDL directory as it contains non-thrift file " + file.Name())
		}
		waitGroup.Add(1)

		go func(fileName, clusterName string) {
			err := s.AddService(fileName, clusterName)
			if err != nil {
				log.Fatal(err)
			}
			waitGroup.Done()
		}(file.Name(), file.Name()[0:len(file.Name())-7])
	}
	waitGroup.Wait()
}

type Store struct {
	mutex sync.RWMutex

	proxyStateListeners    []EventListener
	ServiceMapListeners    []EventListener
	apiStateListeners      []EventListener
	apiValidationListeners []EventListener
	apiRouteListeners      []EventListener
	lbStateListners        []EventListener

	IsOn bool

	ServicesMap map[string]*ServiceMeta

	ProxyAddress          string
	StoreAddress          string
	IdlFolderRelativePath string
	Password              string
}

type ServiceMeta struct {
	ServiceName string
	ClusterName string

	Descriptor *descriptorKeeper
	LbType     string

	APIs map[string]*ApiMeta
}

type ApiMeta struct {
	ServiceName string
	MethodName  string

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
func (s *Store) RegisterApiStateListener(listener EventListener) {
	s.apiStateListeners = append(s.apiStateListeners, listener)
}
func (s *Store) RegisterApiValidationListener(listener EventListener) {
	s.apiValidationListeners = append(s.apiValidationListeners, listener)
}
func (s *Store) RegisterApiRouteListener(listener EventListener) {
	s.apiRouteListeners = append(s.apiRouteListeners, listener)
}
func (s *Store) RegisterLoadBalanceChoiceListener(listener EventListener) {
	s.lbStateListners = append(s.lbStateListners, listener)
}

type EventListener interface {
	OnStatechanged(data ...interface{}) error
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
	meta, ok := s.ServicesMap[serviceName]
	if !ok {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}
	return meta.APIs, nil
}

func (s *Store) AddService(idlFileName, clusterName string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	dk, err := buildDescriptorKeeperFromPath(idlFileName, s.IdlFolderRelativePath)
	if err != nil {
		return err
	}
	sd, err := dk.Get()
	if err != nil {
		return err
	}
	serviceName := sd.Name
	_, ok := s.ServicesMap[serviceName]
	if ok {
		return fmt.Errorf("service %s already exists", serviceName)
	}

	result := map[string]*ApiMeta{}

	for methodName := range sd.Functions {
		route := make(map[string]map[string]bool)
		route["GET"] = make(map[string]bool)
		route["GET"]["/"+serviceName+"/"+methodName] = true
		api := ApiMeta{
			ServiceName:  serviceName,
			MethodName:   methodName,
			ValidationOn: false,
			Routes:       route,
			IsOn:         true,
		}
		result[methodName] = &api
	}

	s.ServicesMap[serviceName] = &ServiceMeta{
		ServiceName: serviceName,
		ClusterName: clusterName,
		Descriptor:  dk,
		//todo
		LbType: "default",
		APIs:   result,
	}
	notifyStatechange(s.ServiceMapListeners, true, s.ServicesMap[serviceName])
	return nil
}

func (s *Store) RemoveService(serviceName string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var err error = nil

	meta, ok := s.ServicesMap[serviceName]
	if !ok {
		err = fmt.Errorf("service %s not found", serviceName)
		return err
	}

	delete(s.ServicesMap, serviceName)
	notifyStatechange(s.ServiceMapListeners, false, meta)
	return err
}

func (s *Store) UpdateService(serviceName, idlFileName, clusterName string) error {
	err := s.RemoveService(serviceName)
	if err != nil {
		return err
	}
	return s.AddService(idlFileName, clusterName)
}

func (s *Store) TurnOnService(serviceName string) error {
	for methodName := range s.ServicesMap[serviceName].APIs {
		err := s.TurnOnAPI(serviceName, methodName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) TurnOffService(serviceName string) error {
	for methodName := range s.ServicesMap[serviceName].APIs {
		err := s.TurnOffAPI(serviceName, methodName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) GetServiceInfo(serviceName string) (*ServiceMeta, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	meta, ok := s.ServicesMap[serviceName]
	if !ok {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}
	return meta, nil
}

func (s *Store) CheckAPIStatus(serviceName, methodName string) (*ApiMeta, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	meta, ok := s.ServicesMap[serviceName]
	if !ok {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}
	apiMeta, ok := meta.APIs[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s not found", methodName)
	}
	return apiMeta, nil
}

func (s *Store) TurnOnAPI(serviceName, methodName string) error {
	api, err := s.CheckAPIStatus(serviceName, methodName)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	api.IsOn = true
	notifyStatechange(s.apiStateListeners, serviceName, methodName, true)
	return nil
}

func (s *Store) TurnOffAPI(serviceName, methodName string) error {
	api, err := s.CheckAPIStatus(serviceName, methodName)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	api.IsOn = false
	notifyStatechange(s.apiStateListeners, serviceName, methodName, false)
	return nil
}

func (s *Store) TurnOnValidation(serviceName, methodName string) error {
	api, err := s.CheckAPIStatus(serviceName, methodName)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	api.ValidationOn = true
	notifyStatechange(s.apiValidationListeners, serviceName, methodName, true)
	return nil
}

func (s *Store) TurnOffValidation(serviceName, methodName string) error {
	api, err := s.CheckAPIStatus(serviceName, methodName)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	api.ValidationOn = false
	notifyStatechange(s.apiValidationListeners, serviceName, methodName, false)
	return nil
}

func (s *Store) AddRoute(serviceName, methodName, httpMethod, url string) error {
	api, err := s.CheckAPIStatus(serviceName, methodName)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	m, ok := api.Routes[httpMethod]
	if !ok {
		m = make(map[string]bool)
		api.Routes[httpMethod] = m
	}
	_, ok = m[url]
	if ok {
		return fmt.Errorf("route %s %s already exists", httpMethod, url)
	}
	m[url] = true
	notifyStatechange(s.apiRouteListeners, serviceName, methodName, url, httpMethod, true)
	return nil
}

func (s *Store) RemoveRoute(serviceName, methodName, httpMethod, url string) error {
	api, err := s.CheckAPIStatus(serviceName, methodName)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	m, ok := api.Routes[httpMethod]
	if !ok {
		m = make(map[string]bool)
		api.Routes[httpMethod] = m
	}
	_, ok = m[url]
	if !ok {
		return fmt.Errorf("route %s %s does not exists", httpMethod, url)
	}
	delete(s.ServicesMap[serviceName].APIs[methodName].Routes[httpMethod], url)
	notifyStatechange(s.apiRouteListeners, serviceName, methodName, url, httpMethod, false)
	return nil
}

func (s *Store) ModifyRoute(serviceName, methodName, httpMethod, url, newMethod, newUrl string) error {
	err := s.RemoveRoute(serviceName, methodName, httpMethod, url)
	if err != nil {
		return err
	}
	return s.AddRoute(serviceName, methodName, newMethod, newUrl)
}

func (s *Store) GetRoutes(serviceName, methodName string) (map[string]map[string]bool, error) {
	api, err := s.CheckAPIStatus(serviceName, methodName)
	if err != nil {
		return nil, err
	}
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return api.Routes, nil
}

func (s *Store) GetLbType(serviceName string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	meta, ok := s.ServicesMap[serviceName]
	if !ok {
		return "", fmt.Errorf("service %s not found", serviceName)
	}
	return meta.LbType, nil
}

func (s *Store) SetLbType(serviceName, lbType string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	meta, ok := s.ServicesMap[serviceName]
	if !ok {
		return fmt.Errorf("service %s not found", serviceName)
	}
	meta.LbType = lbType
	notifyStatechange(s.lbStateListners, meta)
	return nil
}
