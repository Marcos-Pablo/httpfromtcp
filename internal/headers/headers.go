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

var specialChars map[rune]struct{} = map[rune]struct{}{
	'!':  {},
	'#':  {},
	'$':  {},
	'%':  {},
	'&':  {},
	'\'': {},
	'*':  {},
	'+':  {},
	'-':  {},
	'^':  {},
	'_':  {},
	'`':  {},
	'|':  {},
	'~':  {},
}

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

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("poorly formatted header name: %s", key)
	}

	key = strings.TrimSpace(key)
	if !validTokens(key) {
		return 0, false, fmt.Errorf("poorly formatted header name: %s", key)
	}

	value := parts[1]
	value = strings.TrimSpace(value)

	h.Set(key, value)

	return bytesProcessed, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func validTokens(data string) bool {
	for _, char := range data {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			continue
		}

		if char >= '0' && char <= '9' {
			continue
		}

		if _, ok := specialChars[char]; ok {
			continue
		}

		return false
	}

	return len(data) >= 1
}
