package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)

	if err != nil {
		log.Fatalf("couldn't open %s: %s", inputFilePath, err)
	}

	defer file.Close()
	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	line := make([]byte, 0, 64)
	buff := make([]byte, 8)
	for {
		n, err := file.Read(buff)

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
			fmt.Printf("read: %s\n", string(line))
			line = line[:0]
			chunk = chunk[i+1:]
		}
	}

	if len(line) > 0 {
		fmt.Printf("read: %s\n", string(line))
	}
}
