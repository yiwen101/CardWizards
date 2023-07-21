package store

type Store struct {
	Services     map[string]serviceMeta
	proxyOn      bool
	proxyAddress string
	storeAddress string
}

type api struct {
	methodName   string
	validationOn bool
	url          string
}

type serviceMeta struct {
	serviceName string
	idlFileName string
	idlFilePath string
	lbType      string
	isOn        bool
	methods     map[string]api
}

var MetaStore Store = Store{
	Services:     make(map[string]serviceMeta),
	proxyOn:      false,
	proxyAddress: "",
	storeAddress: "",
}
