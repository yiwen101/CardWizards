package utils

import "net/http"

// assign variable "HttpMethods" to a slice of strings representing the http methods

const (
	PkgToIDL           = "../../test/idl"
	TestRPCServerToIDL = "../../../idl"
)

var httpMethods = []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead, http.MethodOptions}

func HTTPMethods() []string {
	return httpMethods
}
