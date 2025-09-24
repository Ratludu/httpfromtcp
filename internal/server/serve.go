package server

import (
	"bytes"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/ratludu/httpfromtcp/internal/request"
	"github.com/ratludu/httpfromtcp/internal/response"
)

type Server struct {
	Closed   atomic.Bool
	Listener net.Listener
	Handler  Handler
	Port     int
}

type Handler func(w *response.Writer, r *request.Request)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func Serve(port int, handlerFunc Handler) (*Server, error) {
	// creates a net.listener and returns a new Server
	// starts listening for requests using a go routine

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	s := &Server{
		Closed:   atomic.Bool{},
		Listener: l,
		Handler:  handlerFunc,
		Port:     port,
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	// Closes the lisenter and server
	s.Closed.Store(true)
	err := s.Listener.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) listen() {
	//  Uses a loop to .Accept new connections as they come in, and handles each one in a new goroutine. I used an atomic.Bool to track whether the server is closed or not so that I can ignore connection errors after the server is closed.

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		if s.Closed.Load() {
			return
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {

	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	buf := new(bytes.Buffer)

	w := response.NewWriter(conn)

	s.Handler(w, req)

	defaultHeaders := response.GetDefaultHeaders(buf.Len(), "text/html")
	statusOk := response.Ok
	err = w.WriteStatusLine(statusOk)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = w.WriteHeaders(defaultHeaders)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, err = w.WriteBody(buf.Bytes())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

}
