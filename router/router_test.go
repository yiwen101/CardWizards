package router

import (
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func TestRouter(t *testing.T) {
	// todo, how to do unit testing
	r := NewRouteManager()
	h := server.Default(
		server.WithHostPorts("127.0.0.1:8080"),
	)
	r.RegisterRoutes(h)
}
