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
	Usersconn = make(map[string]net.Conn)
	maxUsers = 10
)

type Message struct {
	Username string
	Data     string
}


func LaunchServer(port string) {
	
	ClearHistory()
	listner, err := net.Listen("tcp", ":"+port)
	if err != nil { log.Fatal("Error Listen ", err)}
	
	defer listner.Close()

	fmt.Println("Server started on port", port)

	messageChannel := make(chan Message)
	JoinChannel := make(chan string)
	LeaveChannel := make(chan string)

	go handleMessages(messageChannel, JoinChannel, LeaveChannel)
	for {
		conn, err := listner.Accept()
		if err != nil {
			continue
		}
		if len(Usersconn) >= maxUsers {
			fmt.Fprintln(conn, "Server full. Try again later.")
			conn.Close()
			continue
		}
		go handleClient(conn, messageChannel, JoinChannel, LeaveChannel)
	}
}


func handleMessages(messageChannel chan Message, JoinChannel chan string, LeaveChannel chan string) {
	file, err := os.OpenFile("history.txt", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil { return }
	
	defer file.Close()

	for {
		select {
		case msg := <-messageChannel:
			timesending := time.Now().Format("2006-01-02 15:04:05")
			SendMsg := fmt.Sprintf("[%s][%s]: %s", timesending, msg.Username, msg.Data)
			
			file.WriteString(SendMsg + "\n")

			for user, conn := range Usersconn {
				if user != msg.Username {
					fmt.Fprintln(conn, SendMsg)
				}
			}

		case Clientname := <-JoinChannel:
			for user, conn := range Usersconn {
				if user != Clientname {
					fmt.Fprintln(conn, "\n" + fmt.Sprintf("%s has joined the chat.", Clientname))
				}
			}

		case Clientname := <-LeaveChannel:
			for user, conn := range Usersconn {
				if user != Clientname {
					fmt.Fprintln(conn, "\n" + fmt.Sprintf("%s has left the chat.", Clientname))
				}
			}
		}
	}
}



func handleClient(conn net.Conn, messageChannel chan Message, JoinChannel chan string, LeaveChannel chan string) {

	reader := bufio.NewReader(conn)

	conn.Write([]byte(NetCatHeader()))

	username := getValidUsername(conn, reader)
	if username == "" { return }

	Usersconn[username] = conn
	JoinChannel <- username

	sendHistoryToUser(conn)

	for {
		
		timesending := time.Now().Format("2006-01-02 15:04:05")
		prompt := fmt.Sprintf("[%s][%s]:", timesending,username)
		fmt.Fprint(conn, prompt)

		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msg = strings.TrimSpace(msg)

		if isValidMsg(msg){
			messageChannel <- Message{Username: username, Data: msg}
		}
	}

	delete(Usersconn, username)
	LeaveChannel <- username
	
	defer conn.Close()
}
