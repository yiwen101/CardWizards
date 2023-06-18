package temp

import (
	"net/http"
)

const (
	RelativePathToIDL = "../IDL/"
	RefaultHttpMethod = http.MethodPost
	RefaultRoute      = "/:serviceName/:methodName/*furtherRoutes"
)

func Load() {
	buildDescriptorManager()
	buildGenericClients()
}
