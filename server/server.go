package server

import (
	"net"
	"net/url"
)

type Server interface {
	IsAlive() bool
	SetAlive(bool)
	GetUrl() *url.URL
	GetActiveConnections() int
	Serve(listener net.Listener) error
}
