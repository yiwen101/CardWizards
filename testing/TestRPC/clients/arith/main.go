package main

import (
	"context"
	"log"

	arithmatic "github.com/yiwen101/CardWizards/TestRPC/kitex_gen/arithmatic"
	calculator "github.com/yiwen101/CardWizards/TestRPC/kitex_gen/arithmatic/arithmatic"
)

func main() {
	/*
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

		client2, err := noteservice.NewClient("demonote")
		if err != nil {
			log.Fatal(err)
		}
		req2 := &demonote.CreateNoteRequest{}
		resp2, err := client2.CreateNote(context.Background(), req2)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp2)
	*/

	client3, err := calculator.NewClient("arithmatic")
	if err != nil {
		log.Fatal(err)
	}
	req3 := &arithmatic.Request{}
	// set the first arguement of req3 to 1, second arguement to 2
	req3.FirstArguement = 1
	req3.SecondArguement = 2
	resp3, err := client3.Add(context.Background(), req3)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp3)
}
