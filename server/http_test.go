package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransport struct {
	mock.Mock
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("OK")),
		Header:     make(http.Header),
	}, nil
}

func TestGetActiveConnections(t *testing.T) {
	url, _ := url.Parse("http://localhost")
	mockTransport := new(MockTransport)
	server := NewHttpServer(url, &httputil.ReverseProxy{
		Transport: mockTransport,
		Director: func(req *http.Request) {
			// Set the URL to which the proxy will forward requests
			req.URL = url
			req.Header.Set("X-Forwarded-Host", req.Host)
		},
	})

	activeConnections := server.GetActiveConnections()
	assert.Equal(t, 0, activeConnections)

	mockTransport.On("RoundTrip", mock.Anything).Return(&http.Response{StatusCode: http.StatusOK}, nil)

	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	go func() {
		for range 10 {
			server.handler(rr, req)
		}
	}()

	activeConnections = server.GetActiveConnections()
	assert.LessOrEqual(t, activeConnections, 10)

	time.Sleep(100 * time.Millisecond)

	activeConnections = server.GetActiveConnections()
	assert.Equal(t, 0, activeConnections)
}

func TestSetAliveAndIsAlive(t *testing.T) {
	url, _ := url.Parse("http://localhost")
	mockTransport := new(MockTransport)
	server := NewHttpServer(url, &httputil.ReverseProxy{
		Transport: mockTransport,
		Director: func(req *http.Request) {
			// Set the URL to which the proxy will forward requests
			req.URL = url
			req.Header.Set("X-Forwarded-Host", req.Host)
		},
	})

	assert.True(t, server.IsAlive())

	server.SetAlive(false)
	assert.False(t, server.IsAlive())

	server.SetAlive(true)
	assert.True(t, server.IsAlive())
}

func TestServe(t *testing.T) {
	url, _ := url.Parse("http://localhost")
	mockTransport := new(MockTransport)
	server := NewHttpServer(url, &httputil.ReverseProxy{
		Transport: mockTransport,
		Director: func(req *http.Request) {
			// Set the URL to which the proxy will forward requests
			req.URL = url
			req.Header.Set("X-Forwarded-Host", req.Host)
		},
	})

	mockTransport.On("RoundTrip", mock.Anything).Return(&http.Response{StatusCode: http.StatusOK}, nil)

	listener := httptest.NewServer(http.NewServeMux())
	defer listener.Close()

	go func() {
		err := server.Serve(listener.Listener)
		assert.Nil(t, err) // Ensure there is no error
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get(listener.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
