package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/ratludu/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Ok StatusCode = iota
	BadRequest
	InternalServerError
)

func (s StatusCode) GetCode() int {

	switch s {
	case Ok:
		return 200
	case BadRequest:
		return 400
	default:
		return 500
	}
}

func (s StatusCode) GetMessage() string {

	switch s {
	case Ok:
		return "OK"
	case BadRequest:
		return "Bad Request"
	default:
		return "Internal Server Error"
	}

}

func (s StatusCode) CreateHTTPMessage() string {
	return fmt.Sprintf("HTTP/1.1 %d %s", s.GetCode(), s.GetMessage())
}

func GetDefaultHeaders(contentLen int) headers.Headers {

	header := headers.NewHeaders()
	strContentLen := strconv.Itoa(contentLen)
	header["content-length"] = strContentLen
	header["connection"] = "close"
	header["content-type"] = "text/plain"

	return header
}

func (s *StatusCode) WriteHeaders(w io.Writer, headers headers.Headers) error {
	msgHeaders := s.CreateHTTPMessage() + "\r\n"
	for k, v := range headers {
		msgHeaders += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msgHeaders += "\r\n"

	_, err := w.Write([]byte(msgHeaders))
	if err != nil {
		return err
	}

	return nil
}

func WriteBody(w io.Writer, body []byte) error {
	_, err := w.Write([]byte(body))
	if err != nil {
		return err
	}
	return nil
}
