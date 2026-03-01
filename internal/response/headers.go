package response

import (
	"strconv"

	"github.com/Marcos-Pablo/httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	defaultHeaders := headers.NewHeaders()
	defaultHeaders.Set("Connection", "close")
	defaultHeaders.Set("Content-Type", "text/plain")
	defaultHeaders.Set("Content-Length", strconv.Itoa(contentLen))
	return defaultHeaders
}
