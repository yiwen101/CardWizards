package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yiwen101/CardWizards/admin"
	"github.com/yiwen101/CardWizards/pkg/store"
)

// interface; bindValidator/clients to handler; abstraction; object oriented; single duty, unit testing, configuration, singleManager, singleDatabase, use of ok, err and fatal, common.http as const; // complex feature: use the annotation of the thrift file
// intermediate feature: give route at command line; avoid regenerate twice;

func main() {
	// todo, set up logs and tracer?
	// check proposal, feedbasktemplate, milestone2 sample and software engineering/design pattern books
	// add two listeners and produce string information

	store.InfoStore.Load("127.0.0.1:8080", "default", "../../IDL")

	h := server.Default(
		server.WithHostPorts("127.0.0.1:8080"),
	)

	admin.Register(h)
	h.Spin()
}
