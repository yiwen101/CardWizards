package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"log"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/genericserver"
	"github.com/kitex-contrib/registry-nacos/registry"
	//arithmatic "github.com/yiwen101/CardWizards/TestRPC/kitex_gen/arithmatic"
	//calculator "github.com/yiwen101/CardWizards/TestRPC/kitex_gen/arithmatic/arithmatic"
)

// CalculatorImpl implements the last service interface defined in the IDL.
type CalculatorImpl struct{}

/*
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
*/

type GenericServiceImpl struct {
}

type requestStruct struct {
	FirstArguement  int `json:"firstArguement"`
	SecondArguement int `json:"secondArguement"`
}

type responseStruct struct {
	FirstArguement  int `json:"firstArguement"`
	SecondArguement int `json:"secondArguement"`
	Result          int `json:"result"`
}

func (g *GenericServiceImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	str, err := strconv.Unquote(string(jsonBytes))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("strconv.Unquote(string(jsonBytes)) is: %s", str)
	fmt.Println("strconv.Unquote(string(jsonBytes)) is: %s", str)
	fmt.Println("strconv.Unquote(string(jsonBytes)) is: %s", str)
	fmt.Println("strconv.Unquote(string(jsonBytes)) is: %s", str)
	fmt.Println("strconv.Unquote(string(jsonBytes)) is: %s", str)

	var req requestStruct
	err = json.Unmarshal([]byte(str), &req)
	if err != nil {
		return nil, err
	}
	fmt.Println("req is: %v", req)
	fmt.Println("req is: %v", req)
	fmt.Println("req is: %v", req)
	fmt.Println("req is: %v", req)
	fmt.Println("req is: %v", req)
	var resp responseStruct
	resp.FirstArguement = 10
	resp.SecondArguement = 7
	resp.Result = 17
	fmt.Println("resp is: %v", resp)
	fmt.Println("resp is: %v", resp)
	fmt.Println("resp is: %v", resp)
	fmt.Println("resp is: %v", resp)
	fmt.Println("resp is: %v", resp)
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	result := string(respBytes)
	fmt.Println("result is: %s", result)
	fmt.Println("result is: %s", result)
	fmt.Println("result is: %s", result)
	fmt.Println("result is: %s", result)
	fmt.Println("result is: %s", result)
	return string(respBytes), nil
	//return "{\"SecondArguement\":7,\"result\":17,\"firstArguement\":10}", nil

}

func main() {

	p, err := generic.NewThriftFileProvider("arithmatic.thrift", "../../../idl")
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
