package main

import (
	"github.com/Tether-Payments/micadb/simulator"
	"os"
	"strconv"
)

func main() {
	runCount := 1000
	if len(os.Args) == 3 {
		runCount, _ = strconv.Atoi(os.Args[2])
	}

	switch os.Args[1] {
	case "concurrency":
		simulator.Concurrency(runCount)
	case "read":
		simulator.Read()
	case "write":
		simulator.Write(runCount)
	}

}
