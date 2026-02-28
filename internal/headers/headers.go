package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))

	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return len(crlf), true, nil
	}

	bytesProcessed := idx + len(crlf)

	headerStr := string(data[:idx])

	parts := strings.SplitN(headerStr, ":", 2)

	if len(parts) != 2 {
		return bytesProcessed, false, fmt.Errorf("poorly formatted header: %s", headerStr)
	}

	key := parts[0]
	value := parts[1]

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("poorly formatted header name: %s", key)
	}

	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)

	h[key] = value

	return bytesProcessed, false, nil
}
