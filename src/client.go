package src

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"securebox/lib"
)

type Client struct {
}

func (c *Client) Connect() error {
	address := fmt.Sprintf("127.0.0.1:8089")
	serverAddress, _ := net.ResolveTCPAddr("tcp", address)

	conn, err := net.DialTCP("tcp", nil, serverAddress)
	if err != nil {
		return errors.New("Not able to connect to the server\n")
	}
	defer conn.Close()

	fmt.Printf("Connection on %s\n", address)

	sharedKey , privateKey:= lib.Handshake(conn)
	secureConnection := lib.SecureConnection{Conn: conn, SharedKey: sharedKey, PrivateKey:privateKey}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter your message ::")
		// Read up to the newline character
		msg, _ := reader.ReadBytes(0xA)
		// Kill the newline char
		msg = msg[:len(msg)-1]

		_, err := secureConnection.Write(msg)

		response := make([]byte, 1024)

		_, err = secureConnection.Read(response)
		if err != nil {
			fmt.Print("Connection to the server was closed.\n")
			break
		}

		fmt.Printf("%s\n", response)
	}

	return nil
}



