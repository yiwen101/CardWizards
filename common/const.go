package common

import "net/http"

// assign variable "HttpMethods" to a slice of strings representing the http methods

const (
	RelativePathToIDLFromTest  = "../../testing/idl"
	RelativePathToIDLFromTest2 = "../testing/idl"
	GenericPath1               = "/:serviceName/:methodName/*furtherRoutes"
	GenericPath2               = "/*route"
)

var httpMethods = []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead, http.MethodOptions}

func HTTPMethods() []string {
	return httpMethods
}
