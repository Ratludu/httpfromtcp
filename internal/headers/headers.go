package headers

import (
	"bytes"
	"fmt"
	"slices"
	"unicode"
)

const crlf = "\r\n"

var SpecialCharacters = []string{"!", "#", "$", "%", "&", "'", "*", "+", "-", ".", "^", "_", "`", "|", "~"}

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

	if ok := bytes.ContainsRune(splitHeader[0], ' '); ok {
		return 0, false, fmt.Errorf("First part contains spaces")
	}

	cleanSplitValue := bytes.Trim(splitHeader[1], " ")

	if ok := ValidateCharacters(splitHeader[0]); !ok {
		return 0, false, fmt.Errorf("Contains invalid character")
	}

	if len(splitHeader[0]) <= 0 {
		return 0, false, fmt.Errorf("length is not at least one")
	}

	key := bytes.ToLower(splitHeader[0])

	if v, ok := h[string(key)]; ok {
		h[string(key)] = fmt.Sprintf("%s, %s", v, string(cleanSplitValue))
	} else {
		h[string(key)] = string(cleanSplitValue)
	}

	return idx + len(crlf), false, nil
}

func ValidateCharacters(data []byte) bool {
	for _, letter := range data {
		if ok := ValidateChar(letter); !ok {
			return false
		}
	}
	return true
}

func isSpecial(b byte) bool {
	return slices.Contains(SpecialCharacters, string(b))
}

func ValidateChar(b byte) bool {
	return unicode.IsLetter(rune(b)) || unicode.IsDigit(rune(b)) || isSpecial(b)
}
