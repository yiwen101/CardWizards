package store

type Store interface {
	//api, bind, cluster, server, route, plugin
	PutAPI(api *API) (error)
	RemoveAPI(serviceName, methodName string) error
	GetAPI(serviceName, methodName string) (*API, error)
	GetAllAPIs() ([]*API, error)
}

type API struct {
	isOn bool
	validationOn bool
	idlPath string
}

