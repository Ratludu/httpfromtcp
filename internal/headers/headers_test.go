package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeader(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, "localhost:42069", headers["Host"])
	require.NoError(t, err)
	assert.Equal(t, 37, n)
	assert.False(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, done)
}

func TestParse_ValidSingleHeader(t *testing.T) {
	h := NewHeaders()
	n, done, err := h.Parse([]byte("Host: localhost:42069\r\n"))
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, 23, n)
	assert.Equal(t, "localhost:42069", h["Host"])
}

func TestParse_ValidSingleHeaderExtraWhitespace(t *testing.T) {
	h := NewHeaders()
	n, done, err := h.Parse([]byte("     Host:      localhost:42069     \r\n"))
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, "localhost:42069", h["Host"])
	assert.Equal(t, len("     Host:      localhost:42069     \r\n"), n)
}

func TestParse_ValidHeader_NoSpaceAfterColon(t *testing.T) {
	h := NewHeaders()
	n, done, err := h.Parse([]byte("Host:localhost\r\n"))
	require.Error(t, err)
	assert.False(t, done)
	assert.Equal(t, 0, n)
}

func TestParse_InvalidSpaceBeforeColon(t *testing.T) {
	h := NewHeaders()
	n, done, err := h.Parse([]byte("Host : localhost\r\n"))
	require.Error(t, err)
	assert.False(t, done)
	assert.Equal(t, 0, n)
}

func TestParse_Invalid_NoColon(t *testing.T) {
	h := NewHeaders()
	n, done, err := h.Parse([]byte("Host localhost\r\n"))
	require.Error(t, err)
	assert.False(t, done)
	assert.Equal(t, 0, n)
}

func TestParse_Invalid_SpacesInName(t *testing.T) {
	h := NewHeaders()
	n, done, err := h.Parse([]byte("X Bad: value\r\n"))
	require.Error(t, err)
	assert.False(t, done)
	assert.Equal(t, 0, n)
}

func TestParse_ValidMultipleCalls(t *testing.T) {
	h := NewHeaders()
	// First header
	n1, done1, err1 := h.Parse([]byte("Host: localhost\r\n"))
	require.NoError(t, err1)
	assert.False(t, done1)
	assert.Equal(t, "localhost", h["Host"])
	assert.Equal(t, len("Host: localhost\r\n"), n1)

	// Second header
	n2, done2, err2 := h.Parse([]byte("User-Agent: boots\r\n"))
	require.NoError(t, err2)
	assert.False(t, done2)
	assert.Equal(t, "boots", h["User-Agent"])
	assert.Equal(t, len("User-Agent: boots\r\n"), n2)
}

func TestParse_ValidWithExistingHeaders(t *testing.T) {
	h := NewHeaders()
	h["Existing"] = "keep"
	n, done, err := h.Parse([]byte("Content-Type: text/plain\r\n"))
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, "keep", h["Existing"])
	assert.Equal(t, "text/plain", h["Content-Type"])
	assert.Equal(t, len("Content-Type: text/plain\r\n"), n)
}

func TestParse_DoneWhenBlankLine(t *testing.T) {
	h := NewHeaders()
	n, done, err := h.Parse([]byte("\r\n"))
	require.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, 0, n) // consume none per spec for done-at-start
}

func TestParse_TrimValueManySpaces(t *testing.T) {
	h := NewHeaders()
	n, done, err := h.Parse([]byte("Accept:        text/html      \r\n"))
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, "text/html", h["Accept"])
	assert.Equal(t, len("Accept:        text/html      \r\n"), n)
}

func TestParse_EmptyValueAllowed(t *testing.T) {
	h := NewHeaders()
	n, done, err := h.Parse([]byte("X-Flag:\r\n"))
	require.Error(t, err)
	assert.False(t, done)
	assert.Equal(t, 0, n)
}

func TestParse_FieldNameTrimOuterButNoInnerSpaces(t *testing.T) {
	h := NewHeaders()
	n, done, err := h.Parse([]byte("   X-Key: value   \r\n"))
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, "value", h["X-Key"])
	assert.Equal(t, len("   X-Key: value   \r\n"), n)
}
