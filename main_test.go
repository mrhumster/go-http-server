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

func Test_getHello(t *testing.T) {
	responseBody := "Hello, Billy!"
	payload := strings.NewReader("myName=Billy")
	req := httptest.NewRequest(http.MethodPost, "/hello", payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), "serverAddr", "httptest")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	getHello(w, req)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, responseBody, strings.TrimSpace(string(data)))
	req = httptest.NewRequest(http.MethodPost, "/hello", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx = context.WithValue(req.Context(), "serverAddr", "httptest")
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()
	getHello(w, req)
	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "myName", w.Header().Get("x-missing-field"))
}

func Test_getHealth(t *testing.T) {
	responseBody := "✅"
	method := http.MethodGet
	url := "/health"
	req := httptest.NewRequest(method, url, nil)
	ctx := context.WithValue(req.Context(), "serverAddr", "httptest")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	getHealth(w, req)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, responseBody, strings.TrimSpace(string(data)))
}

func Test_getRoot(t *testing.T) {
	responseBody := "This is my website!"
	payload := strings.NewReader("field1=value1")
	req := httptest.NewRequest(http.MethodPost, "/", payload)
	req.Header.Add("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "serverAddr", "httptest")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	getRoot(w, req)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, responseBody, strings.TrimSpace(string(data)))
}

func TestRouting(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", getHello)
	mux.HandleFunc("/health", getHealth)
	mux.HandleFunc("/", getRoot)

	testCases := []struct {
		path     string
		expected int
	}{
		{"/hello", http.StatusBadRequest},
		{"/health", http.StatusOK},
		{"/nonexistent", http.StatusOK},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest("GET", tc.path, nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != tc.expected {
			t.Errorf("path %s: expected status %d, got %d",
				tc.path, tc.expected, rr.Code)
		}
	}
}
