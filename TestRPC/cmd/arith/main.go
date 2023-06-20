package main

import (
	"context"

	"log"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/genericserver"
	"github.com/kitex-contrib/registry-nacos/registry"
	arithmatic "github.com/yiwen101/CardWizards/TestRPC/kitex_gen/arithmatic"
	//calculator "github.com/yiwen101/CardWizards/TestRPC/kitex_gen/arithmatic/arithmatic"
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

type GenericServiceImpl struct {
}

func (g *GenericServiceImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	/*
		var req arithmatic.Request
		jsStr, ok := request.(string)
		if !ok {
			return nil, errors.New("request is not string")
		}

		err = json.Unmarshal([]byte(jsStr), &req)

		if err != nil {
			return nil, err
		}

		response = &arithmatic.Response{FirstArguement: req.FirstArguement, SecondArguement: req.SecondArguement, Result_: req.FirstArguement + req.SecondArguement}
		responseString, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		return responseString, nil
		// use jsoniter or other json parse sdk to assert request
	*/
	return "{\"SecondArguement\":7,\"result\":17,\"firstArguement\":10}", nil

}

func main() {

	p, err := generic.NewThriftFileProvider("arithmatic.thrift", "../../../IDL")
	if err != nil {
		panic(err)
	}

	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}

	r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		panic(err)
	}

	svc := genericserver.NewServer(
		new(GenericServiceImpl),
		g,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "arithmatic"}),
		server.WithRegistry(r))
	if err != nil {
		panic(err)
	}

	/*
		svr := calculator.NewServer(
			new(CalculatorImpl),

		)
	*/

	err = svc.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
