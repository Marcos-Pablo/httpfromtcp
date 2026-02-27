package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(data)

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return nil, fmt.Errorf("could not find CRLF in request-line")
	}

	requestLineText := string(data[:idx])
	parts := strings.Split(requestLineText, " ")

	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", requestLineText)
	}

	method := parts[0]
	target := parts[1]
	httpVersion := parts[2]

	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	if httpVersion != "HTTP/1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpVersion)
	}

	version := strings.Split(httpVersion, "/")[1]

	return &RequestLine{
		HttpVersion:   version,
		Method:        method,
		RequestTarget: target,
	}, nil
}
