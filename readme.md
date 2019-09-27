# Securebox

This is an implementation of securely sending messages between client and server using public/private key mechanism.
Whenever a client sends a message to server it will first do a handshake to exchange the public keys between them. After exchanging the keys the client will encrypt the message using the public key of the server and server will decrypt it using it's private key.  

### Run

Start the server:
```./main -option="server"```

Start the client:
```./main -option="client"```

Enter the message in the shell and the server will reply back the same message

