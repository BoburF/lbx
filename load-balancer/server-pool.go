package loadbalancer

import (
	"errors"
	"net/url"
	"sync"

	"github.com/BoburF/lbx/server"
)

type ServerPool struct {
	pool    []server.Server
	mu      sync.RWMutex
	current int
}

type ServerPoolUnit struct {
	Url   *url.URL
	Force int
}

func NewServerPool(urls []ServerPoolUnit) ServerPool {
	pool := make([]server.Server, 0)

	for _, v := range urls {
		srv := server.NewHttpServer(v.Url, v.Force)
		for range v.Force {
			pool = append(pool, srv)
		}
	}

	return ServerPool{
		pool:    pool,
		current: 0,
	}
}

func (sp *ServerPool) Next() (server.Server, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	n := len(sp.pool)
	if n == 0 {
		return nil, errors.New("pool is empty")
	}

	start := sp.current
	for {
		srv := sp.pool[sp.current]
		sp.current = (sp.current + 1) % n

		if srv.IsAlive() {
			return srv, nil
		}

		if sp.current == start {
			return nil, errors.New("no server is alive")
		}
	}
}
