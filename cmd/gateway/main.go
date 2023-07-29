package main

import (
	"flag"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/pprof"
	"github.com/yiwen101/CardWizards/pkg/admin"
	"github.com/yiwen101/CardWizards/pkg/proxy"
	"github.com/yiwen101/CardWizards/pkg/store"
)

var (
	addr          = flag.String("addr", "127.0.0.1:8080", "Addr: http request entrypoint")
	pathIDL       = flag.String("idl", "IDL", "Path: idl file path")
	adminPassword = flag.String("pwd", "", "addr Store Password")
	pprofOn       = flag.Bool("pprof", false, "Enable pprof")
	//addrPPROF                     = flag.String("addr-pprof", "", "Addr: pprof addr")
	//limitBytesCachingMB           = flag.Uint64("limit-caching", 64, "Limit(MB): MB for caching size")
	//version                       = flag.Bool("version", false, "Show version info")
)

// todo, set up logs and tracer?
// todo, add pprof information to frontEnd
// extension: add config.Option, eg: 在较大 request size 下（request size > 1M），推荐使用 go net 网络库加流式。在其他场景下，推荐使用 netpoll 网络库，会获得极致的性能。

func main() {
	flag.Parse()

	store.InfoStore.Load(*addr, *pathIDL, *adminPassword)

	h := server.Default(
		server.WithHostPorts(*addr),
	)

	h.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // Add the allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Add the allowed request headers
		ExposeHeaders:    []string{"Content-Length"},                          // Expose additional response headers if needed
		AllowCredentials: true,                                                // Allow credentials (e.g., cookies, authorization headers)
		MaxAge:           12 * time.Hour,                                      // Set the preflight request cache duration
	}))

	admin.Register(h)
	if *pprofOn {
		pprof.Register(h)
	}
	proxy.Register(h)
	pprof.Register(h)
	h.Spin()
}
