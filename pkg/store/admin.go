package store

type admin interface {
	GetAllServiceNames() (map[string]ServiceMeta, error)
	CheckProxyStatus() (bool, error)
	TurnOnProxy() error  //proxyGate
	TurnOffProxy() error //proxyGate

	GetAPIs(serviceName string) (map[string]ApiMeta, error)
	AddService(serviceName, idlFileName, clusterName string) error    //validator, client, descriptor, router
	UpdateService(serviceName, idlFileName, clusterName string) error //validator, client, descriptor, router
	RemoveService(serviceName string) error                           //validator, client, descriptor, router
	TurnOnService(serviceName string) error
	TurnOffService(serviceName string) error
	GetServiceInfo(serviceName string) (ServiceMeta, error)

	CheckAPIStatus(serviceName, methodName string) (ApiMeta, error)
	TurnOnAPI(serviceName, methodName string) error                 //APIGate
	TurnOffAPI(serviceName, methodName string) error                //APIGate
	TurnOnValidation(serviceName, methodName string) error          //validator
	TurnOffValidation(serviceName, methodName string) error         //validator
	AddRoute(serviceName, methodName, url, httpMethod string) error //router
	ModifyRoute(url, httpMethod, newUrl, newMethod string) error    //router
	RemoveRoute(url, httpMethod string) error                       //router
	GetRoutes(serviceName, methodName string) (map[string]map[string]bool, error)
}
