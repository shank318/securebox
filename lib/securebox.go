/*
A simple wrapper to send and receive messages over the secure connection
pass the clients public key and own private key to it
and call read and write to send/receive data.
 */

package lib

import (
	"bytes"
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/nacl/box"
	"net"
)

type BoxMessage struct {
	msg   []byte
	nonce [24]byte
}

func (s *BoxMessage) toByteArray() []byte {
	return append(s.nonce[:], s.msg[:]...)
}

func CreateBoxMessage(sm []byte) BoxMessage {
	var nonce [24]byte

	//Create a random nonce for each sent message
	// To provide more security
	nonceArray := sm[:24]
	copy(nonce[:], nonceArray)

	// Trim out all unnecessary bytes
	msg := bytes.Trim(sm[24:], "\x00")

	return BoxMessage{msg: msg, nonce: nonce}
}

type SecureConnection struct {
	Conn       *net.TCPConn
	SharedKey  *[32]byte
	PrivateKey *[32]byte
}

func (s *SecureConnection) Read(p []byte) (int, error) {
	message := make([]byte, 2048)

	n, err := s.Conn.Read(message)

	secureMessage := CreateBoxMessage(message)
	decryptedMessage, ok := box.Open(nil, secureMessage.msg, &secureMessage.nonce, s.SharedKey, s.PrivateKey)

	if !ok {
		return 0, errors.New("Unable to decrypt the message.\n")
	}

	n = copy(p, decryptedMessage)

	return n, err
}

func (s *SecureConnection) Write(p []byte) (int, error) {
	var nonce [24]byte

	//Create a random nonce for each sent message
	// To provide more security
	rand.Read(nonce[:])

	encryptedMessage := box.Seal(nil, p, &nonce, s.SharedKey, s.PrivateKey)
	sm := BoxMessage{msg: encryptedMessage, nonce: nonce}

	return s.Conn.Write(sm.toByteArray())
}

/**
 Handshake will generate a random public/private key
and exchange public keys by
Reading the public key coming from the client
and writing it's public key to the client
 */
func Handshake(conn *net.TCPConn) (*[32] byte, *[32] byte) {
	var peerKey [32]byte

	//Generate key pair
	publicKey, privateKey, _ := box.GenerateKey(rand.Reader)

	// Sending the public key to the client vice-a-versa
	conn.Write(publicKey[:])

	//Receive the public key to the server vice-a-versa
	peerKeyArray := make([]byte, 32)
	conn.Read(peerKeyArray)

	copy(peerKey[:], peerKeyArray)

	return &peerKey, privateKey
}
