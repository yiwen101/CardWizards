package main

import (
	"context"
	"log"

	"github.com/cloudwego/kitex/client"
	"github.com/yiwen101/CardWizards/kitex_gen/demouser"
	"github.com/yiwen101/CardWizards/kitex_gen/demouser/userservice"
)

func main() {
	client, err := userservice.NewClient("demonuser", client.WithHostPorts("0.0.0.0:8888"))
	if err != nil {
		log.Fatal(err)
	}
	req := &demouser.CreateUserRequest{}
	resp, err := client.CreateUser(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}
