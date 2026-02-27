package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const crlf = "\r\n"
const bufferSize = 8

type ParserState int

const (
	Initialized ParserState = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	State       ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		State: Initialized,
	}

	for req.State != Done {
		if readToIndex >= len(buff) {
			tmp := make([]byte, len(buff)*2)
			copy(tmp, buff)
			buff = tmp
		}

		bytesRead, err := reader.Read(buff[readToIndex:])

		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.State != Done {
					return nil, fmt.Errorf("incomplete request")
				}
			}
			return nil, err
		}

		readToIndex += bytesRead

		bytesParsed, err := req.parse(buff[:readToIndex])

		if err != nil {
			return nil, err
		}

		copy(buff, buff[bytesParsed:readToIndex])
		readToIndex -= bytesParsed
	}

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case Initialized:
		reqLine, bytesRead, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if bytesRead == 0 {
			return 0, nil
		}

		r.RequestLine = *reqLine
		r.State = Done

		return bytesRead, nil
	case Done:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return nil, 0, nil
	}
	bytesRead := idx + len(crlf)

	requestLineText := string(data[:idx])
	parts := strings.Split(requestLineText, " ")

	if len(parts) != 3 {
		return nil, bytesRead, fmt.Errorf("poorly formatted request-line: %s", requestLineText)
	}

	method := parts[0]
	target := parts[1]
	httpVersion := parts[2]

	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, bytesRead, fmt.Errorf("invalid method: %s", method)
		}
	}

	if httpVersion != "HTTP/1.1" {
		return nil, bytesRead, fmt.Errorf("unrecognized HTTP-version: %s", httpVersion)
	}

	version := strings.Split(httpVersion, "/")[1]

	return &RequestLine{
		HttpVersion:   version,
		Method:        method,
		RequestTarget: target,
	}, bytesRead, nil
}
