package configuer

import (
	"context"
	"encoding/json"
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

/*
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
*/

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
}

/*
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
*/

func generalHandler(ctx context.Context, c *app.RequestContext) {
	path := string(c.URI().Path())
	method := string(c.Method())
	c.String(http.StatusOK, "received: "+method+"\n")

	routeManager, err := router.GetRouteManager()
	if err != nil {
		hlog.Fatal("Internal Server Error in getting the route manager: ", err)
	}
	api, ok := routeManager.GetRoute(method, path)
	if !ok {
		c.String(http.StatusBadRequest, "invalid route: "+path+"\n")
		return
	}

	hm, err := service.GetHandlerManager()
	if err != nil {
		hlog.Fatal("Internal Server Error in getting the handler manager: ", err)
	}

	f, err := hm.HandlerForRoute(api.ServiceName, api.MethodName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error in getting the handler: "+err.Error())
		return
	}
	f(ctx, c)
}

func Register(r *server.Hertz) {
	// todo: other http methods

	r.Any("/*path", generalHandler)

	type update struct {
		ServiceName string
		FileName    string
		Dir         string
	}
	r.PUT("/update", func(ctx context.Context, c *app.RequestContext) {
		jsonbytes, err := c.Body()
		if err != nil {
			c.String(http.StatusBadRequest, "invalid body: "+err.Error())
			return
		}
		jsonStr := string(jsonbytes)
		c.String(http.StatusOK, "received: "+jsonStr)

		update := update{}

		err = json.Unmarshal(jsonbytes, &update)
		if err != nil {
			c.String(http.StatusBadRequest, "unmarshalError: "+err.Error())
			return
		}

		err = clients.ClientManager.UpdateClient(update.ServiceName, update.FileName, update.Dir)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error in updating the client: "+err.Error())
			return
		}
		c.String(http.StatusOK, "updated")
	})

	log.Println("routes and handlers registered to the server")
}

var useFolder = flag.String("useFolder", "./IDL", "relative path to IDL folder")

func Load() {
	flag.Parse()
	descriptor.BuildDescriptorManager(*useFolder)
	clients.BuildGenericClients(*useFolder)

	routeManager, err := router.GetRouteManager()
	if err != nil {
		hlog.Fatal("Internal Server Error in getting the route manager: ", err)
	}
	err = routeManager.InitRoute()
	if err != nil {
		hlog.Fatal("Internal Server Error in getting the route manager: ", err)
	}

}

// Option, customised routes? annotated routes?
// must, absolute file path to the idl file
