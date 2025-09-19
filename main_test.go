package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHello(t *testing.T) {
	responseBody := "Hello, Billy!"
	payload := strings.NewReader("myName=Billy")

	req := httptest.NewRequest(http.MethodPost, "/hello", payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), "serverAddr", "httptest")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	// Run Handler
	getHello(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, responseBody, strings.TrimSpace(string(data)))

}
