package main

import (
	arithmatic "arithmatic/kitex_gen/arithmatic/arithmatic"
	"log"
)

func main() {
	svr := arithmatic.NewServer(new(ArithmaticImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
