package main

import (
	arithmatic "arithmatic/kitex_gen/arithmatic"
	"context"
)

// ArithmaticImpl implements the last service interface defined in the IDL.
type ArithmaticImpl struct{}

// Add implements the ArithmaticImpl interface.
func (s *ArithmaticImpl) Add(ctx context.Context, req *arithmatic.Request) (resp *arithmatic.Response, err error) {
	// TODO: Your code here...
	return &arithmatic.Response{Result_: req.FirstArguement + req.SecondArguement}, nil
}

// Sub implements the ArithmaticImpl interface.
func (s *ArithmaticImpl) Sub(ctx context.Context, req *arithmatic.Request) (resp *arithmatic.Response, err error) {
	// TODO: Your code here...
	return &arithmatic.Response{Result_: req.FirstArguement - req.SecondArguement}, nil
}

// Mul implements the ArithmaticImpl interface.
func (s *ArithmaticImpl) Mul(ctx context.Context, req *arithmatic.Request) (resp *arithmatic.Response, err error) {
	// TODO: Your code here...
	return &arithmatic.Response{Result_: req.FirstArguement * req.SecondArguement}, nil
}

// Div implements the ArithmaticImpl interface.
func (s *ArithmaticImpl) Div(ctx context.Context, req *arithmatic.Request) (resp *arithmatic.Response, err error) {
	// TODO: Your code here...
	return &arithmatic.Response{Result_: req.FirstArguement / req.SecondArguement}, nil
}
