package main

import (
	"log"
	"os"
	"strconv"

	"github.com/Stroevik3/DistrArithmeticExpresCalc/internal/app/agent"
)

func main() {
	cfg := agent.NewConfig()
	sCp := os.Getenv("COMPUTING_POWER")
	if sCp != "" {
		cp, err := strconv.Atoi(sCp)
		if err == nil {
			cfg.ComputingPower = cp
		}
	}

	if err := agent.Start(cfg); err != nil {
		log.Fatal(err)
	}
}
