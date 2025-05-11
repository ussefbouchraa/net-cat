package netcat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	users    = make(map[string]net.Conn)
	maxUsers = 10
)

type Message struct {
	Username string
	Data     string
}

func Clear() {
	file, err := os.OpenFile("history.txt", os.O_TRUNC, 0o644)
	if err != nil {
		log.Fatal("Error clear file", err)
		return
	}
	file.Close()
}

func AcceptConnections(port string) {
	Clear()
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Error Listen", err)
	}
	defer ln.Close()

	fmt.Println("Server started on port", port)

	messageChannel := make(chan Message)
	JoinChannel := make(chan string)
	LeaveChannel := make(chan string)

	go handleMessages(messageChannel, JoinChannel, LeaveChannel)
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		if len(users) >= maxUsers {
			fmt.Fprintln(conn, "Server full. Try again later.")
			conn.Close()
			continue
		}
		go handleClient(conn, messageChannel, JoinChannel, LeaveChannel)
	}
}

func IsvalidClientName(username string) bool {
	if len(username) == 0 {
		return false
	}
	for _, ch := range username {
		if !(ch >= 'a' && ch <= 'z') && !(ch >= 'A' && ch <= 'Z') {
			return false
		}
	}
	return true
}

func handleMessages(messageChannel chan Message, JoinChannel chan string, LeaveChannel chan string) {
	file, err := os.OpenFile("history.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal("Error Open File", err)
		return
	}
	defer file.Close()
	for {
		select {
		case msg := <-messageChannel:
			timesending := time.Now().Format("2006-01-02 15:04:05")
			SendMsg := fmt.Sprintf("[%s][%s]: %s", timesending, msg.Username, msg.Data)

			file.WriteString(SendMsg + "\n")

			for user, conn := range users {
				if user != msg.Username {
					fmt.Fprintln(conn, SendMsg)
				}
			}

		case Clientname := <-JoinChannel:
			for user, conn := range users {
				if user != Clientname {
					fmt.Fprintln(conn, fmt.Sprintf("%s has joined the chat.", Clientname))
				}
			}

		case Clientname := <-LeaveChannel:
			for user, conn := range users {
				if user != Clientname {
					fmt.Fprintln(conn, fmt.Sprintf("%s has left the chat.", Clientname))
				}
			}
		}
	}
}

func handleClient(conn net.Conn, messageChannel chan Message, JoinChannel chan string, LeaveChannel chan string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	
	netcat := `Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\\dS"qML
 |    '.       | '' \\Zq
_)      \\.___.,|     .'
\\____   )MMMMMP|   .'
     '-'       '--'
`

	fmt.Fprintln(conn, netcat)

	var username string

	for {
		fmt.Fprintln(conn, "Enter your Name: ")
		rawUsername, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading username:", err)
			return
		}

		username = strings.TrimSpace(rawUsername)
		if !IsvalidClientName(username) {
			fmt.Fprintln(conn, "Invalid username. Try again.")
			continue
		}

		if _, exists := users[username]; exists {
			fmt.Fprintln(conn, "Username already taken. Try again.")
			continue
		}

		users[username] = conn
		JoinChannel <- username
		break
	}

	Data, err := os.ReadFile("history.txt")
	if err != nil {
		log.Fatal("Error open file history.txt")
		return
	}

	for _, v := range Data {
		fmt.Fprint(conn, string(v))
	}

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error Reding Message", err)
			break
		}
		msg = strings.TrimSpace(msg)
		if len(msg) > 0 {
			messageChannel <- Message{Username: username, Data: msg}
		}
	}

	delete(users, username)
	LeaveChannel <- username
}
