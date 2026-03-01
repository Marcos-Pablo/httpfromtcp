package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/Marcos-Pablo/httpfromtcp/internal/request"
	"github.com/Marcos-Pablo/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
	handler  Handler
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(w, headers)
	w.Write(messageBytes)
}

func Serve(port int, handler Handler) (*Server, error) {
	portStr := strconv.Itoa(port)
	listener, err := net.Listen("tcp", ":"+portStr)

	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
		isClosed: atomic.Bool{},
		handler:  handler,
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
	req, err := request.RequestFromReader(conn)

	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}

	buff := new(bytes.Buffer)
	hErr := s.handler(buff, req)

	if hErr != nil {
		hErr.Write(conn)
		return
	}

	err = response.WriteStatusLine(conn, response.StatusOK)

	if err != nil {
		log.Printf("Error writing status line to response: %s", err)
		return
	}

	defaultHeaders := response.GetDefaultHeaders(buff.Len())
	err = response.WriteHeaders(conn, defaultHeaders)

	if err != nil {
		log.Printf("Error writing headers to response: %s", err)
		return
	}

	_, err = conn.Write(buff.Bytes())

	if err != nil {
		log.Printf("Error writing body to response: %s", err)
		return
	}
}
