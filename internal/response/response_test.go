package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {

	// Test Ok
	s := Ok
	msg := s.CreateHTTPMessage()
	assert.Equal(t, "HTTP/1.1 200 OK", msg)

	// Test bad request
	s = BadRequest
	msg = s.CreateHTTPMessage()
	assert.Equal(t, "HTTP/1.1 400 Bad Request", msg)

	// Test bad request
	s = InternalServerError
	msg = s.CreateHTTPMessage()
	assert.Equal(t, "HTTP/1.1 500 Internal Server Error", msg)
}
