package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os/exec"
)

// not much needed as of now
type ServerOpts struct {
	Address    string
	MaxSession int
	// these are fine for now ig
}

// not much needed as of now
type Server struct {
	CurrSessions   []string
	MaxSession     int
	RequestingAddr string
}

// Starts the server and creates in-memory sessions
func StartServer(opts *ServerOpts) {
	conn, err := net.Dial("tcp", opts.Address)
	if err != nil {
		log.Fatalf("Failed to connect to the SSH-F Server: %v\n", err)
	}
	defer conn.Close()
	conn.Write([]byte("ok\n"))
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to establish connection: %v\n", err)
	}
	if line != "ok" {
		fmt.Println("Recieved Reply from ssh-f server: ", line)
		log.Fatalf("Failed to establish connection with ssh-f server.\n")
	}

	go handleConn(conn)

}

func handleConn(conn net.Conn) {

}
