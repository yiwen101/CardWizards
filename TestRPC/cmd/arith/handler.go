package main

import (
	"context"

	arithmatic "github.com/yiwen101/TiktokXOrbital-CardWizards/kitex_gen/arithmatic"
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
