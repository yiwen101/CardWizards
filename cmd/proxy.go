package main

import (
	"flag"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/yiwen101/CardWizards/admin"
	"github.com/yiwen101/CardWizards/pkg/store/descriptor"
	"github.com/yiwen101/CardWizards/router"
	"github.com/yiwen101/CardWizards/service/clients"
)

// interface; bindValidator/clients to handler; abstraction; object oriented; single duty, unit testing, configuration, singleManager, singleDatabase, use of ok, err and fatal, common.http as const; // complex feature: use the annotation of the thrift file
// intermediate feature: give route at command line; avoid regenerate twice;

var (
	ProxyAddress   = flag.String("addr", "127.0.0.1:8080", "proxy address")
	IDLFolederPath = flag.String("idl", "././IDL", "idl folder path")
)

func main() {
	// todo, set up logs and tracer?
	// check proposal, feedbasktemplate, milestone2 sample and software engineering/design pattern books
	// add two listeners and produce string information

	admin.Load()

	h := server.Default(
		server.WithHostPorts("127.0.0.1:8080"),
	)

	admin.Register(h)
	h.Spin()
}

func Load() {
	flag.Parse()
	descriptor.BuildDescriptorManager(*IDLFolederPath)
	clients.BuildGenericClients(*IDLFolederPath)

	routeManager, err := router.GetRouteManager()
	if err != nil {
		hlog.Fatal("Internal Server Error in getting the route manager: ", err)
	}
	err = routeManager.InitRoute()
	if err != nil {
		hlog.Fatal("Internal Server Error in getting the route manager: ", err)
	}

}
