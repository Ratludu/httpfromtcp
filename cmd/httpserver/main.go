package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/ratludu/httpfromtcp/internal/headers"
	"github.com/ratludu/httpfromtcp/internal/request"
	"github.com/ratludu/httpfromtcp/internal/response"
	"github.com/ratludu/httpfromtcp/internal/server"
)

const port = 42069

func main() {

	HandlerFunc := func(w *response.Writer, r *request.Request) {

		switch {
		case r.RequestLine.RequestTarget == "/yourproblem":
			w.WriteStatusLine(response.BadRequest)
			message, n := response400()
			defaultHeaders := response.GetDefaultHeaders(n, "text/html")
			w.WriteHeaders(defaultHeaders)
			w.WriteBody(message)
		case r.RequestLine.RequestTarget == "/myproblem":
			w.WriteStatusLine(response.InternalServerError)
			message, n := response500()
			defaultHeaders := response.GetDefaultHeaders(n, "text/html")
			w.WriteHeaders(defaultHeaders)
			w.WriteBody(message)
		case strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin"):
			resp, err := http.Get("https://httpbin.org" + r.RequestLine.RequestTarget[8:])
			if err != nil {
				fmt.Println("Error:", err)
				w.WriteStatusLine(response.InternalServerError)
				defaultHeaders := response.GetDefaultHeaders(0, "text/plain")
				w.WriteHeaders(defaultHeaders)
				return
			}
			w.WriteStatusLine(response.Ok)
			defaultHeaders := response.GetDefaultHeaders(0, "text/plain")
			defaultHeaders.Del("content-length")
			defaultHeaders.Set("transfer-encoding", "chunked")
			defaultHeaders.Set("trailer", "X-Content-SHA256, X-Content-Length")
			w.WriteHeaders(defaultHeaders)

			fullBody := []byte("")
			for {
				data := make([]byte, 32)
				n, err := resp.Body.Read(data)
				if err != nil {
					break
				}

				fullBody = append(fullBody, data[:n]...)
				w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
				w.WriteBody(data[:n])
				w.WriteBody([]byte("\r\n"))
			}
			w.WriteBody([]byte("0\r\n"))
			trailingHeader := headers.NewHeaders()
			sum := fmt.Sprintf("%x", sha256.Sum256(fullBody))
			sumLength := strconv.Itoa(len(fullBody))

			trailingHeader.Set("X-Content-SHA256", sum)
			trailingHeader.Set("X-Content-Length", sumLength)
			w.WriteHeaders(trailingHeader)
		case strings.HasPrefix(r.RequestLine.RequestTarget, "/video"):
			vid, err := os.ReadFile("assets/vim.mp4")
			if err != nil {
				fmt.Println("Error:", err)
				w.WriteStatusLine(response.InternalServerError)
				defaultHeaders := response.GetDefaultHeaders(0, "text/plain")
				w.WriteHeaders(defaultHeaders)
				return
			}
			w.WriteStatusLine(response.Ok)
			defaultHeaders := response.GetDefaultHeaders(len(vid), "video/mp4")
			w.WriteHeaders(defaultHeaders)
			w.WriteBody(vid)
		default:
			w.WriteStatusLine(response.Ok)
			message, n := response200()
			defaultHeaders := response.GetDefaultHeaders(n, "text/html")
			w.WriteHeaders(defaultHeaders)
			w.WriteBody(message)
		}
	}
	s, err := server.Serve(port, HandlerFunc)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer s.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func response200() ([]byte, int) {

	message := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

	return []byte(message), len(message)

}

func response400() ([]byte, int) {

	message := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

	return []byte(message), len(message)

}

func response500() ([]byte, int) {

	message := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`

	return []byte(message), len(message)

}
