package main

import (
	"net/http"
	"snippetbox.samedarslan28.net/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {
	// Create a new instance of our application struct. For now, this just
	// contains a couple of mock loggers (which discard anything written to
	// them).
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, body := ts.get(t, "/ping")

	assert.Equal(t, statusCode, http.StatusOK)

	assert.Equal(t, body, "OK")
}
