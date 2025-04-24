package server

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type httpServer struct {
	url          *url.URL
	alive        bool
	mux          sync.RWMutex
	connections  int
	reverseProxy *httputil.ReverseProxy
}

func (s *httpServer) GetActiveConnections() int {
	s.mux.RLock()
	connections := s.connections
	s.mux.RUnlock()
	return connections
}

func (s *httpServer) SetAlive(alive bool) {
	s.mux.Lock()
	s.alive = alive
	s.mux.Unlock()
}

func (s *httpServer) IsAlive() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.alive
}

func (s *httpServer) GetUrl() *url.URL {
	return s.url
}

func (s *httpServer) Serve(listener net.Listener) error {
	srv := &http.Server{
		Handler: http.HandlerFunc(s.handler),
	}

	return srv.Serve(listener)
}

func (s *httpServer) handler(w http.ResponseWriter, r *http.Request) {
	s.mux.Lock()
	s.connections++
	s.mux.Unlock()

	defer func() {
		s.mux.Lock()
		s.connections--
		s.mux.Unlock()
	}()

	s.reverseProxy.ServeHTTP(w, r)
}

func NewHttpServer(url *url.URL, rp *httputil.ReverseProxy) httpServer {
	return httpServer{
		url:          url,
		alive:        true,
		reverseProxy: rp,
	}
}
