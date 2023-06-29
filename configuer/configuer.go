package configuer

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/yiwen101/CardWizards/common"
	"github.com/yiwen101/CardWizards/common/descriptor"
	"github.com/yiwen101/CardWizards/router"
	"github.com/yiwen101/CardWizards/service"
	"github.com/yiwen101/CardWizards/service/clients"
)

type toRegist struct {
	httpMethod string
	path       string
	handler    func(ctx context.Context, c *app.RequestContext)
}

func generateToRegists() error {
	if tRs != nil {
		return nil
	}
	ls1, err := generateGenericToRegists()
	if err != nil {
		log.Fatal("Internal Server Error in getting the handler: ", err)
		return err
	}
	ls2, err := generateRegularToRegists()
	if err != nil {
		log.Fatal("Internal Server Error in getting the handler: ", err)
		return err
	}

	tRs = append(ls1, ls2...)
	return nil
}

func generateGenericToRegists() ([]toRegist, error) {
	hm, err := service.GetHandlerManager()
	if err != nil {
		hlog.Fatal("Internal Server Error in getting the handler manager: ", err)
	}

	ls := make([]toRegist, 7)

	// how to write this as const? since arrays are not constants

	for i, method := range common.HTTPMethods() {
		handlerFunc, err := hm.HandlerForAnnotatedRoutes(method)
		if err != nil {
			log.Fatal("Internal Server Error in getting the handler: ", err)
			return nil, err
		}

		tR := toRegist{
			httpMethod: method,
			path:       common.GenericPath2,
			handler:    handlerFunc,
		}
		ls[i] = tR
	}
	return ls, nil
	// todo
	// complex feature: use the annotation of the thrift file
	// intermediate feature: give route at command line
}

func generateRegularToRegists() ([]toRegist, error) {
	hm, err := service.GetHandlerManager()
	if err != nil {
		hlog.Fatal("Internal Server Error in getting the handler manager: ", err)
	}

	rm, err := router.GetRouteManager()
	if err != nil {
		hlog.Fatal("Internal Server Error in getting the route manager: ", err)
	}

	routes, err := rm.BuildAndGetRoute()
	if err != nil {
		hlog.Fatal("Internal Server Error in getting the routes: ", err)
	}

	ls := make([]toRegist, len(routes))

	for i, route := range routes {
		http, path := route.GetRoute()
		handler, err := hm.HandlerForRoute(route.ServiceName, route.MethodName)
		if err != nil {
			log.Fatal("Internal Server Error in getting the handler: ", err)
			return nil, err
		}
		toRegist := toRegist{
			httpMethod: http,
			path:       path,
			handler:    handler,
		}
		ls[i] = toRegist
	}
	return ls, nil
}

var tRs []toRegist

func Register(r *server.Hertz) {
	generateToRegists()

	for _, tR := range tRs {
		switch tR.httpMethod {
		case http.MethodGet:
			r.GET(tR.path, tR.handler)
			continue
		case http.MethodPost:
			r.POST(tR.path, tR.handler)
			continue
		case http.MethodPut:
			r.PUT(tR.path, tR.handler)
			continue
		case http.MethodDelete:
			r.DELETE(tR.path, tR.handler)
			continue
		case http.MethodPatch:
			r.PATCH(tR.path, tR.handler)
			continue
		case http.MethodHead:
			r.HEAD(tR.path, tR.handler)
			continue
		case http.MethodOptions:
			r.OPTIONS(tR.path, tR.handler)
			continue
		default:
			log.Println(" unsupported http method, invalid route ")
			continue
		}
	}
	log.Println("routes and handlers registered to the server")
}

var useFolder = flag.String("useFolder", "./IDL", "relative path to IDL folder")

func Load() {
	flag.Parse()
	descriptor.BuildDescriptorManager(*useFolder)
	clients.BuildGenericClients(*useFolder)
}

// Option, customised routes? annotated routes?
// must, absolute file path to the idl file
