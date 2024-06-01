package main

import (
	"log"

	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/orchestratorserv"
)

func main() {
	cfg := orchestratorserv.NewConfig()

	if err := orchestratorserv.Start(cfg); err != nil {
		log.Fatal(err)
	}
}
