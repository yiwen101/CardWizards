package main

import (
	"flag"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/pprof"
	"github.com/yiwen101/CardWizards/pkg/admin"
	"github.com/yiwen101/CardWizards/pkg/proxy"
	"github.com/yiwen101/CardWizards/pkg/store"
)

// interface; bindValidator/clients to handler; abstraction; object oriented; single duty, unit testing, configuration, singleManager, singleDatabase, use of ok, err and fatal, common.http as const; // complex feature: use the annotation of the thrift file
// intermediate feature: give route at command line; avoid regenerate twice;
var (
	addr          = flag.String("addr", "127.0.0.1:8080", "Addr: http request entrypoint")
	pathIDL       = flag.String("idl", "../IDL", "Path: idl file path")
	adminPassword = flag.String("addr-store-pwd", "", "addr Store Password")
	//addrPPROF                     = flag.String("addr-pprof", "", "Addr: pprof addr")

	//limitBytesCachingMB           = flag.Uint64("limit-caching", 64, "Limit(MB): MB for caching size")
	//version                       = flag.Bool("version", false, "Show version info")
)

func main() {
	flag.Parse()
	// todo, set up logs and tracer?
	// check proposal, feedbasktemplate, milestone2 sample and software engineering/design pattern books
	// add two listeners and produce string information

	store.InfoStore.Load(*addr, *pathIDL, *adminPassword)

	h := server.Default(
		server.WithHostPorts(*addr),
	)

	admin.Register(h)
	proxy.Register(h)
	pprof.Register(h)
	// what is the cache solution for apigateway with json body?
	// extension: add config.Option
	// 在较大 request size 下（request size > 1M），推荐使用 go net 网络库加流式。在其他场景下，推荐使用 netpoll 网络库，会获得极致的性能。
	/*When using Service Registration and Discovery, Spin will register the service into a registry center when starting up, and use signalWaiter to monitor service exceptions. Only by using Spin can we support graceful shutdown.*/

	h.Spin()
}
