package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/ratludu/httpfromtcp/internal/headers"
)

const crlf = "\r\n"

const bufferSize = 8

const (
	initialized state = iota
	done
	requestStateParsingHeaders
	requestStateParsingBody
)

type state int

type Request struct {
	RequestLine RequestLine
	State       state
	Headers     headers.Headers
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func newRequest() *Request {
	return &Request{
		State:   initialized,
		Headers: headers.NewHeaders(),
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	request := newRequest()
	buf := make([]byte, bufferSize)
	readToIndex := 0
	for !(request.State == done) {

		if len(buf) == readToIndex {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf[:readToIndex])
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])

		readToIndex += n
		readN, perr := request.parse(buf[:readToIndex])
		if perr != nil {
			return nil, perr
		}

		copy(buf, buf[readN:readToIndex])
		readToIndex -= readN

		if err == io.EOF {
			if request.State != done {
				return nil, fmt.Errorf("incomplete request at EOF")
			}
			if readN == 0 {
				break
			}
			continue
		}

		if err != nil {
			return nil, err
		}
		if n == 0 && readN == 0 {
			break
		}
	}

	return request, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case initialized:
		n, rl, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *rl
		r.State = requestStateParsingHeaders
		return n, nil
	case done:
		return 0, fmt.Errorf("parser is done")
	case requestStateParsingHeaders:
		n, isDone, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if isDone {
			r.State = requestStateParsingBody
		}

		return n, nil
	case requestStateParsingBody:

		val, err := r.Headers.Get("content-length")
		if err != nil {
			r.State = done
			return len(data), nil
		}

		remaining := val - len(r.Body)
		if remaining < 0 {
			r.State = done
			return 0, nil
		}
		consumed := min(remaining, len(data))
		r.Body = append(r.Body, data[:consumed]...)
		if len(r.Body) > val {
			return consumed, fmt.Errorf("Error: body is longer than content-length, body: %d, content-length: %d", len(r.Body), val)
		}
		if len(r.Body) == val {
			r.State = done
			return consumed, nil
		}

		return consumed, nil

	default:
		return 0, fmt.Errorf("unknown state")
	}
}
func (r *Request) parse(data []byte) (int, error) {

	totalBytesParsed := 0
	for r.State != done {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return totalBytesParsed, nil
		}
		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

func parseRequestLine(requestLine []byte) (int, *RequestLine, error) {

	idx := bytes.Index(requestLine, []byte(crlf))
	if idx == -1 {
		return 0, nil, nil
	}

	read := idx + len(crlf)

	request, err := requestLineFromString(string(requestLine[:idx]))
	if err != nil {
		return 0, nil, err
	}

	return read, request, nil
}

func requestLineFromString(line string) (*RequestLine, error) {

	requestLineParts := strings.Split(line, " ")
	if len(requestLineParts) != 3 {
		return &RequestLine{}, fmt.Errorf("Request line does not have 3 parts")
	}

	// check method is upper
	method := requestLineParts[0]
	for _, letter := range method {
		if unicode.IsLower(letter) {
			return &RequestLine{}, fmt.Errorf("Method is not all upper case")
		}
	}

	requestTarget := requestLineParts[1]

	version := strings.Split(requestLineParts[2], "/")
	if len(version) != 2 {
		return &RequestLine{}, fmt.Errorf("HTTP version does not have 2 parts e.g. HTTP/1.1")
	}

	versionHTTP := version[0]
	if versionHTTP != "HTTP" {
		return &RequestLine{}, fmt.Errorf("First part does not equal HTTP e.g. HTTP/1.1")
	}

	versionDigit := version[1]
	if versionDigit != "1.1" {
		return &RequestLine{}, fmt.Errorf("HTTP version is not 1.1")
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionDigit,
	}, nil

}
