package main

import (
	"snippetbox.samedarslan28.net/internal/assert"
	"testing"
)

func TestRoutesNotFound(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, _ := ts.get(t, "/test/invalid")

	assert.Equal(t, statusCode, 404)

}
