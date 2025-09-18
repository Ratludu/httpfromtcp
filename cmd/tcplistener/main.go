package main

import (
	"fmt"
	"github.com/ratludu/httpfromtcp/internal/request"
	"net"
)

type config struct {
	Bytes int
	Addr  string
}

func main() {

	cfg := config{
		Bytes: 8,
		Addr:  ":42069",
	}

	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Server listening on localhost%s\n", cfg.Addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Message accepted")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		fmt.Println("Connection closed")
	}

}
