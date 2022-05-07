package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAuth this function test for middleware auth function
func TestAuth(t *testing.T) {
	testCases := []struct {
		desc           string
		authentication string
		statusCode     int
	}{
		{"Success", "0000", http.StatusOK},
		{"Error", "auth1234", http.StatusUnauthorized},
	}

	for i := range testCases {
		var handle = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

		req := httptest.NewRequest(http.MethodPost, "/car", nil)
		req.Header.Add("authorize", testCases[i].authentication)

		w := httptest.NewRecorder()
		handle(w, req)
		a := Auth(handle)
		a.ServeHTTP(w, req)
		assert.Equal(t, testCases[i].statusCode, w.Code, "Test Case Failed")
	}
}
