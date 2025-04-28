package lb

import "net/url"

type Server struct {
	force   int
	adrress url.URL
}

func (s *Server) GetForce() int {
	return s.force
}

func (s *Server) GetAdress() url.URL {
	return s.adrress
}
