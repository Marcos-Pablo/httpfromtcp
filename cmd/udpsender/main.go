package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const network = "udp"
const address = ":42069"

func main() {
	udpAdd, err := net.ResolveUDPAddr(network, address)

	if err != nil {
		log.Fatalf("could not resolve udp address: %s", err)
	}

	conn, err := net.DialUDP(network, nil, udpAdd)

	if err != nil {
		log.Fatalf("could not dial up udp connection: %s", err)
	}

	defer conn.Close()

	buff := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		message, err := buff.ReadBytes('\n')

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		_, err = conn.Write(message)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			os.Exit(1)
		}
	}
}
