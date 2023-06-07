// Code generated by Kitex v0.5.2. DO NOT EDIT.
package calculator

import (
	arithmatic "github.com/yiwen101/TiktokXOrbital-CardWizards/kitex_gen/arithmatic"
	server "github.com/cloudwego/kitex/server"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler arithmatic.Calculator, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}
