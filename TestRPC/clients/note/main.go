package main

import (
	"context"
	"log"

	"github.com/cloudwego/kitex/client"

	"github.com/yiwen101/TiktokXOrbital-CardWizards/kitex_gen/demonote"
	"github.com/yiwen101/TiktokXOrbital-CardWizards/kitex_gen/demonote/noteservice"
)

func main() {

	client2, err := noteservice.NewClient("demonote", client.WithHostPorts("0.0.0.0:8888"))
	if err != nil {
		log.Fatal(err)
	}
	req2 := &demonote.CreateNoteRequest{}
	resp2, err := client2.CreateNote(context.Background(), req2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp2)

}
