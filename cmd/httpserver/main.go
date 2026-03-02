package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Marcos-Pablo/httpfromtcp/internal/headers"
	"github.com/Marcos-Pablo/httpfromtcp/internal/request"
	"github.com/Marcos-Pablo/httpfromtcp/internal/response"
	"github.com/Marcos-Pablo/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/html") {
		handlerProxyHTML(w, req)
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		handlerProxy(w, req)
		return
	}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handler400(w, req)
	case "/myproblem":
		handler500(w, req)
	default:
		handler200(w, req)
	}
}

func handlerProxyHTML(w *response.Writer, req *request.Request) {
	URL := "https://httpbin.org/html"
	fmt.Println("Proxying to", URL)
	resp, err := http.Get(URL)

	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusOK)
	h := headers.NewHeaders()
	ct := resp.Header.Get("Content-Type")
	h.Set("Content-Type", ct)
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Connection", "close")
	h.Set("Trailer", "X-Content-SHA256, X-Content-Length")
	w.WriteHeaders(h)

	const maxChunkSize = 1024
	chunks := make([]byte, 0, maxChunkSize)
	buff := make([]byte, maxChunkSize)
	for {
		n, err := resp.Body.Read(buff)
		fmt.Printf("%d bytes read\n", n)

		if n > 0 {
			chunks = append(chunks, buff[:n]...)
			_, err = w.WriteChunkedBody(buff[:n])
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}
		}

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Printf("Error reading response from %s - error: %s", URL, err)
			break
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Printf("Error writing last chunk to body: %s", err)
	}

	trailers := headers.NewHeaders()
	hash := sha256.Sum256(chunks)
	hashStr := hex.EncodeToString(hash[:])
	trailers.Set("X-Content-SHA256", hashStr)
	trailers.Set("X-Content-Length", strconv.Itoa(len(chunks)))

	err = w.WriteTrailers(trailers)
	if err != nil {
		log.Printf("Error writing trailers to body: %s", err)
	}
}

func handlerProxy(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
	URL := "https://httpbin.org" + target
	fmt.Println("Proxying to", URL)
	resp, err := http.Get(URL)

	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusOK)
	h := headers.NewHeaders()
	ct := resp.Header.Get("Content-Type")
	h.Set("Content-Type", ct)
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Connection", "close")
	w.WriteHeaders(h)

	const maxChunkSize = 1024
	buff := make([]byte, maxChunkSize)
	for {
		n, err := resp.Body.Read(buff)
		fmt.Printf("%d bytes read\n", n)

		if n > 0 {
			_, err = w.WriteChunkedBody(buff[:n])
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}
		}

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Printf("Error reading response from %s - error: %s", URL, err)
			break
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Printf("Error writing last chunk to body: %s", err)
	}
	_, err = w.WriteBodyDone()
	if err != nil {
		log.Printf("Error writing last chunk to body: %s", err)
	}

}

func handler400(w *response.Writer, _ *request.Request) {
	title := "400 Bad Request"
	h1 := "Bad Request"
	p := "Your request honestly kinda sucked."
	body := hydrateHTML(title, h1, p)

	h := response.GetDefaultHeaders(len(body))
	h.ReplaceOrSet("Content-Type", "text/html")

	w.WriteStatusLine(response.StatusBadRequest)
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))
}

func handler500(w *response.Writer, _ *request.Request) {
	title := "500 Internal Server Error"
	h1 := "Internal Server Error"
	p := "Okay, you know what? This one is on me."
	body := hydrateHTML(title, h1, p)

	h := response.GetDefaultHeaders(len(body))
	h.ReplaceOrSet("Content-Type", "text/html")

	w.WriteStatusLine(response.StatusInternalServerError)
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))
}

func handler200(w *response.Writer, _ *request.Request) {
	title := "200 OK"
	h1 := "Success!"
	p := "Your request was an absolute banger."
	body := hydrateHTML(title, h1, p)

	h := response.GetDefaultHeaders(len(body))
	h.ReplaceOrSet("Content-Type", "text/html")

	w.WriteStatusLine(response.StatusOK)
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))

}

func hydrateHTML(title, h1, p string) string {
	return fmt.Sprintf(`<html>
  <head>
    <title>%s</title>
  </head>
  <body>
    <h1>%s</h1>
    <p>%s</p>
  </body>
</html>`, title, h1, p)
}
