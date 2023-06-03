package main

import (
	"log"
	api "testServer/kitex_gen/api/arithmatic"
)

func main() {
	svr := api.NewServer(new(ArithmaticImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
