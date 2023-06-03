package main

import (
	"context"
	api "testServer/kitex_gen/api"
)

// ArithmaticImpl implements the last service interface defined in the IDL.
type ArithmaticImpl struct{}

// Add implements the ArithmaticImpl interface.
func (s *ArithmaticImpl) Add(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	// TODO: Your code here...
	return &api.Response{Result_: req.FirstArguement + req.SecondArguement}, nil
}

// Sub implements the ArithmaticImpl interface.
func (s *ArithmaticImpl) Sub(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	// TODO: Your code here...
	return &api.Response{Result_: req.FirstArguement - req.SecondArguement}, nil
}

// Mul implements the ArithmaticImpl interface.
func (s *ArithmaticImpl) Mul(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	// TODO: Your code here...
	return &api.Response{Result_: req.FirstArguement * req.SecondArguement}, nil
}

// Div implements the ArithmaticImpl interface.
func (s *ArithmaticImpl) Div(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	// TODO: Your code here...
	return &api.Response{Result_: req.FirstArguement / req.SecondArguement}, nil
}
