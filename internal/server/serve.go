package server

import (
	"fmt"
	"net"
	"sync/atomic"
)

type Server struct {
	Closed   atomic.Bool
	Listener net.Listener
	Port     int
}

func Serve(port int) (*Server, error) {
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
	message := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello World!\n"
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
