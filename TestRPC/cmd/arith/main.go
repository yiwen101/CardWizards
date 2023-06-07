package main

import (
	arithmatic "github.com/yiwen101/TiktokXOrbital-CardWizards/kitex_gen/arithmatic/calculator"
	"log"
)

func main() {
	svr := arithmatic.NewServer(new(CalculatorImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
