package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"snippetbox.samedarslan28.net/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	ping(responseRecorder, request)
	responseResult := responseRecorder.Result()

	assert.Equal(t, responseResult.StatusCode, http.StatusOK)

	defer responseResult.Body.Close()

	body, err := io.ReadAll(responseResult.Body)
	if err != nil {
		t.Fatal(err)
		return
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "Ok")

}
