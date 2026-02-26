package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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
		ch := getLinesChannel(conn)

		for msg := range ch {
			fmt.Println(msg)
		}

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.Reader) <-chan string {
	ch := make(chan string)

	line := make([]byte, 0, 64)
	buff := make([]byte, 8)
	go func() {
		defer close(ch)
		for {
			n, err := f.Read(buff)

			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatalf("couldn't read from file: %s", err)
			}

			chunk := buff[:n]
			for {
				i := bytes.IndexByte(chunk, '\n')
				if i == -1 {
					line = append(line, chunk...)
					break
				}

				line = append(line, chunk[:i]...)
				ch <- string(line)
				line = line[:0]
				chunk = chunk[i+1:]
			}
		}

		if len(line) > 0 {
			ch <- string(line)
		}
	}()

	return ch
}
