package headers

import (
	"bytes"
	"fmt"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return 0, true, nil
	}

	cleanedHeader := bytes.Trim(data[:idx], " ")
	splitHeader := bytes.Split(cleanedHeader, []byte(": "))
	if len(splitHeader) != 2 {
		return 0, false, fmt.Errorf("Length of split header is not 2")
	}

	if i := bytes.ContainsRune(splitHeader[0], ' '); i {
		return 0, false, fmt.Errorf("First part contains spaces")
	}

	cleanSplitValue := bytes.Trim(splitHeader[1], " ")
	h[string(splitHeader[0])] = string(cleanSplitValue)

	return idx + len(crlf), false, nil
}
