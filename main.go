package main

import (
	"flag"
	"fmt"
	"securebox/src"
)

func main() {

	option := flag.String("option", "server", "Set if running the server.")
	flag.Parse()

	if *option == "server" {
		fmt.Printf("Server started...\n")
		s := src.Server{}
		s.Run()
	} else if *option == "client" {
		fmt.Printf("Client started...\n")
		c := src.Client{}
		err := c.Connect()
		if err != nil {
			fmt.Print(err)
		}
	}
}
