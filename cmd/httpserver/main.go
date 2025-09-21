package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ratludu/httpfromtcp/internal/request"
	"github.com/ratludu/httpfromtcp/internal/response"
	"github.com/ratludu/httpfromtcp/internal/server"
)

const port = 42069

func main() {

	HandlerFunc := func(w io.Writer, r *request.Request) *server.HandlerError {

		switch r.RequestLine.RequestTarget {
		case "/yourproblem":
			return &server.HandlerError{
				StatusCode: response.BadRequest,
				Message:    "Your problem is not my problem\n",
			}
		case "/myproblem":
			return &server.HandlerError{
				StatusCode: response.InternalServerError,
				Message:    "Woopsie, my bad\n",
			}

		default:
			_, err := w.Write([]byte("All good, frfr\n"))
			if err != nil {
				return &server.HandlerError{
					StatusCode: response.InternalServerError,
					Message:    "Internal Server Error\n",
				}
			}
		}
		return nil
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
