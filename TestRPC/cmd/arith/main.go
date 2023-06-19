package main

import (
	"context"

	"log"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/registry"
	arithmatic "github.com/yiwen101/CardWizards/TestRPC/kitex_gen/arithmatic"
	calculator "github.com/yiwen101/CardWizards/TestRPC/kitex_gen/arithmatic/arithmatic"
)

// CalculatorImpl implements the last service interface defined in the IDL.
type CalculatorImpl struct{}

// Add implements the CalculatorImpl interface.
func (s *CalculatorImpl) Add(ctx context.Context, request *arithmatic.Request) (resp *arithmatic.Response, err error) {
	// TODO: Your code here...
	return &arithmatic.Response{FirstArguement: request.FirstArguement, SecondArguement: request.SecondArguement, Result_: request.FirstArguement + request.SecondArguement}, nil
}

// Subtract implements the CalculatorImpl interface.
func (s *CalculatorImpl) Subtract(ctx context.Context, request *arithmatic.Request) (resp *arithmatic.Response, err error) {
	// TODO: Your code here...
	return &arithmatic.Response{FirstArguement: request.FirstArguement, SecondArguement: request.SecondArguement, Result_: request.FirstArguement - request.SecondArguement}, nil
}

// Multiply implements the CalculatorImpl interface.
func (s *CalculatorImpl) Multiply(ctx context.Context, request *arithmatic.Request) (resp *arithmatic.Response, err error) {
	// TODO: Your code here...
	return &arithmatic.Response{FirstArguement: request.FirstArguement, SecondArguement: request.SecondArguement, Result_: request.FirstArguement * request.SecondArguement}, nil
}

// Divide implements the CalculatorImpl interface.
func (s *CalculatorImpl) Divide(ctx context.Context, request *arithmatic.Request) (resp *arithmatic.Response, err error) {
	// TODO: Your code here...
	return &arithmatic.Response{FirstArguement: request.FirstArguement, SecondArguement: request.SecondArguement, Result_: request.FirstArguement / request.SecondArguement}, nil
}

func (s *CalculatorImpl) TestValidator(ctx context.Context, request *arithmatic.TestValidator) (resp *arithmatic.Response, err error) {
	// TODO: Your code here...
	return &arithmatic.Response{FirstArguement: 17, SecondArguement: 17, Result_: 17}, nil
}

func main() {

	r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		panic(err)
	}

	svr := calculator.NewServer(
		new(CalculatorImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "arithmatic"}),
		server.WithRegistry(r),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
