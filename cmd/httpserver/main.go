package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		title := "400 Bad Request"
		h1 := "Bad Request"
		p := "Your request honestly kinda sucked."
		body := hydrateHTML(title, h1, p)

		h := response.GetDefaultHeaders(len(body))
		h.ReplaceOrSet("Content-Type", "text/html")

		w.WriteStatusLine(response.StatusBadRequest)
		w.WriteHeaders(h)
		w.WriteBody([]byte(body))
	case "/myproblem":
		title := "500 Internal Server Error"
		h1 := "Internal Server Error"
		p := "Okay, you know what? This one is on me."
		body := hydrateHTML(title, h1, p)

		h := response.GetDefaultHeaders(len(body))
		h.ReplaceOrSet("Content-Type", "text/html")

		w.WriteStatusLine(response.StatusInternalServerError)
		w.WriteHeaders(h)
		w.WriteBody([]byte(body))
	default:
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
