package main

import (
	"fmt"
	"github.com/8bitstout/orderPushScheduler"
	"github.com/8bitstout/orderPushScheduler/srfh"
	"os"
)

const (
	DEFAULT_ADDRESS = "localhost"
	DEFAULT_PORT    = "50051"
)

func main() {
	fmt.Println("Running order push scheduler")
	arguments := os.Args
	command := arguments[1]
	switch command {
	case "server":
		{
			s := orderPushScheduler.NewScheduler()
			fmt.Println("Scheduler running at", DEFAULT_ADDRESS+":"+DEFAULT_PORT)
			s.Run()
		}
	case "client":
		{
			url := fmt.Sprint(DEFAULT_ADDRESS, ":", DEFAULT_PORT)
			c := srfh.MakeClient(url)
			c.Schedule()
		}
	default:
		fmt.Println("Command not recognized")
	}
}
