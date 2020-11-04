package main

import (
	"fmt"
	"github.com/8bitstout/orderPushScheduler/srfh"
	"os"
)

const (
	DEFAULT_ADDRESS = "localhost"
	DEFAULT_PORT    = "8080"
)

func main() {
	fmt.Println("Running order push scheduler")
	arguments := os.Args
	command := arguments[1]
	switch command {
	case "srfh":
		{
			s := srfh.MakeSrfh(DEFAULT_PORT)
			fmt.Println("Server running at", DEFAULT_ADDRESS+":"+DEFAULT_PORT)
			s.Run()
		}
	case "server":
		// run server
	case "client":
		{
			url := fmt.Sprint("ws://", DEFAULT_ADDRESS, ":", DEFAULT_PORT)
			c := srfh.MakeClient(url)
			c.ListenForMessages()
			c.WriteMessage("I connected!")
			for {

			}
		}
	default:
		fmt.Println("Command not recognized")
	}
}
