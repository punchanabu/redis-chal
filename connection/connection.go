package connection

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/command/handler"
	"github.com/codecrafters-io/redis-starter-go/command/parser"
	"github.com/codecrafters-io/redis-starter-go/config"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func HandleConnection(conn net.Conn, store *store.Store, config *config.ReplicaConfig) {
	defer conn.Close()
	fmt.Println("Connection from ", conn.RemoteAddr().String())

	for {
		// Read data from the connection
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Failed to read data: ", err.Error())
			break
		}

		// Handle PING check
		if string(buffer[:n]) == "+PING\r\n" {
			handler.HandleCommand("PING", nil, store, config)
			break
		}

		// Decode the data
		command, argument, err := parser.Decode(buffer[:n])
		if err != nil {
			fmt.Println("Failed to decode data: ", err.Error())
			break
		}

		// Handle the command
		response := handler.HandleCommand(command, argument, store, config)

		// Encode the response
		encodedResponse := parser.Encode(response)

		// Write the response
		conn.Write([]byte(encodedResponse))
	}
}