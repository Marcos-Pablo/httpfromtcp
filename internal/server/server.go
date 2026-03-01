package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	portStr := strconv.Itoa(port)
	listener, err := net.Listen("tcp", ":"+portStr)

	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
		isClosed: atomic.Bool{},
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()

		if s.isClosed.Load() {
			break
		}

		if err != nil {
			log.Printf("Error accepting connection: %s", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	response := []byte(`HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 13

Hello World!
`)
	_, err := conn.Write(response)

	if err != nil {
		log.Printf("Error writing response to tcp connection: %s", err)
	}
}
