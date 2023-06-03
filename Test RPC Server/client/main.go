package main

import (
	"context"
	"log"
	"testServer/kitex_gen/api"
	"testServer/kitex_gen/api/arithmatic"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

func main() {
	c, err := arithmatic.NewClient("example", client.WithHostPorts("0.0.0.0:8888"))
	if err != nil {
		log.Fatal(err)
	}
	// make a request with arguement 1 and 2
	resp, err := c.Add(context.Background(), &api.Request{FirstArguement: 1, SecondArguement: 2}, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}
