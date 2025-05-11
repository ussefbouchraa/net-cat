package main

import (
	"fmt"
	"os"

	FU "netcat/connections"
)

func main() {
	port := "8989"

	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	} else if len(os.Args) == 2 {
		port = os.Args[1]
	}

	FU.AcceptConnections(port)
}
