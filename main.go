package main

import (
	"fmt"
	"os"

	N "netcat/static"
)

func main() {
	port := "8989"

	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	} else if len(os.Args) == 2 {
		port = os.Args[1]
		if port == "" { os.Stderr.WriteString("Err :Empty Port\n") ; return  }
	}
	N.LaunchServer(port)
}
