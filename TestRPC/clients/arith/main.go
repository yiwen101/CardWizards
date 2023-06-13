package main

import (
	"context"
	"log"

	"github.com/cloudwego/kitex/client"
	"github.com/kitex-contrib/registry-nacos/resolver"
	"github.com/yiwen101/TiktokXOrbital-CardWizards/kitex_gen/arithmatic"
	"github.com/yiwen101/TiktokXOrbital-CardWizards/kitex_gen/arithmatic/calculator"
)

func main() {

	r, err := resolver.NewDefaultNacosResolver()
	if err != nil {
		panic(err)
	}

	client3, err := calculator.NewClient(
		"arith",
		client.WithHostPorts("0.0.0.0:8888"),
		client.WithResolver(r),
	)
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
