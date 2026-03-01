package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Marcos-Pablo/httpfromtcp/internal/request"
)

const network = "tcp"
const address = ":42069"

func main() {
	listener, err := net.Listen(network, address)

	if err != nil {
		log.Fatalf("unable to open a %s connection on address %s", network, address)
	}

	defer listener.Close()

	fmt.Printf("Listening for %s traffic on %s\n", network, address)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Unable to accept connection: %s", err)
			continue
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())
		req, err := request.RequestFromReader(conn)

		if err != nil {
			log.Printf("Unable to parse request: %s", err)
			continue
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}

		fmt.Println("Body:")
		fmt.Println(string(req.Body))

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}
