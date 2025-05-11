package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func readFromServer(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Server disconnected.", err)
			return
		}
		fmt.Print(msg)
	}
}

func writeToServer(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text() + "\n"
		conn.Write([]byte(text))
		if text == "exit\n" {
			fmt.Println("👋 Closing chat.")
			return
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8989")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	reader := bufio.NewReader(conn)

	servmessage, _ := reader.ReadString(':')
	fmt.Print(servmessage)

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		name := scanner.Text() + "\n"
		conn.Write([]byte(name))
	}

	go readFromServer(conn)

	writeToServer(conn)
}
