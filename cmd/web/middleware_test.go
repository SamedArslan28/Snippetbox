package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"snippetbox.samedarslan28.net/internal/assert"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	next := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(responseRecorder, request)
	responseResult := responseRecorder.Result()

	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, responseResult.Header.Get("Content-Security-Policy"), expectedValue)
	assert.Equal(t, "origin-when-cross-origin", responseRecorder.Header().Get("Referrer-Policy"))
	assert.Equal(t, "nosniff", responseRecorder.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "deny", responseRecorder.Header().Get("X-Frame-Options"))
	assert.Equal(t, "0", responseRecorder.Header().Get("X-XSS-Protection"))
	assert.Equal(t, http.StatusOK, responseResult.StatusCode)

	defer responseResult.Body.Close()
	body, err := io.ReadAll(responseResult.Body)
	if err != nil {
		t.Fatal(err)
		return
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
