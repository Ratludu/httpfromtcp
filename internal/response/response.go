package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/ratludu/httpfromtcp/internal/headers"
)

type StatusCode int
type WriterState int

const (
	Ok StatusCode = iota
	BadRequest
	InternalServerError
)

const (
	StateStatusLine WriterState = iota
	StateHeaders
	StateBody
)

type Writer struct {
	writer io.Writer
	state  WriterState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		state:  StateStatusLine,
	}
}
func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := statusCode.CreateHTTPMessage()
	_, err := w.writer.Write([]byte(statusLine))
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	var msgHeaders string
	for k, v := range headers {
		msgHeaders += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msgHeaders += "\r\n"

	_, err := w.writer.Write([]byte(msgHeaders))
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write([]byte(p))
	if err != nil {
		return n, err
	}
	return n, nil
}

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
	return fmt.Sprintf("HTTP/1.1 %d %s\n", s.GetCode(), s.GetMessage())
}

func GetDefaultHeaders(contentLen int, contentType string) headers.Headers {

	header := headers.NewHeaders()
	strContentLen := strconv.Itoa(contentLen)
	header["content-length"] = strContentLen
	header["connection"] = "close"
	header["content-type"] = contentType

	return header
}
