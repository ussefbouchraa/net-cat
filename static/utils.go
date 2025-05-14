package netcat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func NetCatHeader() string {
	return `Welcome to TCP-Chat!
          _nnnn_                      
         dGGGGMMb                     
        @p~qp~~qMb                    
        M|@||@) M|                    
        @,----.JM|                    
       JS^\__/  qKL                   
      dZP        qKRb                 
     dZP          qKKb                
    fZP            SMMb               
    HZM            MMMM               
    FqM            MMMM               
__| ".        |\dS"qML                
|    '.       | '' \Zq                
_)      \.___.,|     '.               
\____   )MMMMMP|   .'                 
     '-'       '--'                   
`
}

	

func IsvalidClientName(username string) bool {
	if len(username) == 0  || len(username) > 10 {
		return false
	}
	for _, ch := range username {
		if !(ch >= 'a' && ch <= 'z') && !(ch >= 'A' && ch <= 'Z') {
			return false
		}
	}
	return true
}

func sendHistoryToUser(conn net.Conn) {
	data, err := os.ReadFile("history.txt")
	if err != nil{ 	
		fmt.Fprintln(conn, "No history available.")
		return 
	}
	
	fmt.Fprint(conn, string(data))
}


func ClearHistory() {
	file, err := os.OpenFile("history.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}


func getValidUsername(conn net.Conn, reader *bufio.Reader) string {
	for {
		fmt.Fprintln(conn, "Enter your Name: ")
		rawUsername, err := reader.ReadString('\n')
		if err != nil { return "" }

		username := strings.TrimSpace(rawUsername)

		if !IsvalidClientName(username) {
			fmt.Fprintln(conn, "Invalid username. Try again.")
			continue
		}

		if _, exists := Usersconn[username]; exists {
			fmt.Fprintln(conn, "Username already taken. Try again.")
			continue
		}
		return username
	}
}

