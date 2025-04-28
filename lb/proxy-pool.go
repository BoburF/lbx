package lb

import (
	"errors"
	"sync"
	"time"

	"github.com/BoburF/lbx/proxy"
)

type ProxyPool struct {
	pool             []*proxy.HttpProxy
	current          int
	retryTimeSeconds int
	mu               sync.RWMutex
}

func NewProxyPool(proxies []*proxy.HttpProxy, retryTimeSeconds int) *ProxyPool {
	pp := &ProxyPool{
		pool:             proxies,
		retryTimeSeconds: retryTimeSeconds,
	}

	go pp.pingWithTimeout()
	return pp
}

func (pp *ProxyPool) Next() (*proxy.HttpProxy, error) {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	n := len(pp.pool)
	if n == 0 {
		return nil, errors.New("proxy pool is empty")
	}

	start := pp.current
	for {
		proxy := pp.pool[pp.current]

		pp.current = (pp.current + 1) % n

		if proxy.GetStatusIsAlive() {
			return proxy, nil
		}

		if pp.current == start {
			break
		}
	}

	return nil, errors.New("no alive proxies found")
}

func (pp *ProxyPool) pingWithTimeout() {
	ticker := time.NewTicker(time.Duration(pp.retryTimeSeconds) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pp.mu.Lock()
		for _, p := range pp.pool {
			alive := p.IsAlive()
			p.SetAlive(alive)
		}
		pp.mu.Unlock()
	}
}
