package src

import (
	"fmt"
	"net"
	"os"
	"securebox/lib"
)

type Server struct {
}

func (s *Server) Run() {
	address := fmt.Sprintf("127.0.0.1:8089")
	networkAddress, _ := net.ResolveTCPAddr("tcp", address)

	listener, err := net.ListenTCP("tcp", networkAddress)
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}

	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			fmt.Print(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close()
	sharedKey, privateKey := lib.Handshake(conn)
	secureConnection := lib.SecureConnection{Conn: conn, SharedKey: sharedKey, PrivateKey: privateKey}

	for {
		msg := make([]byte, 1024)
		_, err := secureConnection.Read(msg)

		if err != nil {
			break
		}

		secureConnection.Write(msg)
	}
}
