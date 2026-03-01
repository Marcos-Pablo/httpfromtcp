package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/Marcos-Pablo/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)
const crlf = "\r\n"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var reasonPhrase string
	switch statusCode {
	case StatusOK:
		reasonPhrase = "HTTP/1.1 200 OK"
	case StatusBadRequest:
		reasonPhrase = "HTTP/1.1 400 Bad Request"
	case StatusInternalServerError:
		reasonPhrase = "HTTP/1.1 500 Internal Server Error"
	}

	_, err := w.Write([]byte(reasonPhrase + crlf))

	if err != nil {
		return fmt.Errorf("Error writing status line: %s", err)
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	defaultHeaders := headers.NewHeaders()
	defaultHeaders.Set("Connection", "close")
	defaultHeaders.Set("Content-Type", "text/plain")
	defaultHeaders.Set("Content-Length", strconv.Itoa(contentLen))
	return defaultHeaders
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := fmt.Fprintf(w, "%s: %s%s", k, v, crlf)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, crlf)
	if err != nil {
		return err
	}
	return nil
}
