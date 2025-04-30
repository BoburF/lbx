package lb

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/BoburF/lbx/config"
	"github.com/BoburF/lbx/proxy"
)

type LoadBalancer struct {
	proxyPool ProxyPool
}

func NewLoadBalancer(serverConfigs []config.ServerConfig, retryTimeSecond int) (LoadBalancer, error) {
	servers := make([]Server, 0)

	for _, serverConfig := range serverConfigs {
		parsedUrl, err := url.Parse(serverConfig.Adress)
		if err != nil {
			return LoadBalancer{}, err
		}
		servers = append(servers, Server{adrress: *parsedUrl, force: serverConfig.Force})
	}

	pool := make([]*proxy.HttpProxy, 0)

	for _, server := range servers {
		for range server.GetForce() {
			target := server.GetAdress()
			localTarget := target

			rp := httputil.NewSingleHostReverseProxy(&localTarget)
			rp.Director = func(req *http.Request) {
				req.URL.Scheme = target.Scheme
				req.URL.Host = target.Host
				req.Host = target.Host
			}

			proxy := proxy.HttpProxy{
				Adress: target,
				Proxy:  *rp,
			}
			proxy.SetAlive(true)

			pool = append(pool, &proxy)
		}
	}

	pp := NewProxyPool(pool, retryTimeSecond)
	return LoadBalancer{
		proxyPool: *pp,
	}, nil
}

func (lb *LoadBalancer) Server(host string, port int) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", lb.handle)

	server := &http.Server{
		Addr:         net.JoinHostPort(host, itoa(port)),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Load balancer listening on %s:%d", host, port)
	return server.ListenAndServe()
}

func (lb *LoadBalancer) handle(w http.ResponseWriter, r *http.Request) {
	proxy, err := lb.proxyPool.Next()
	if err != nil {
		http.Error(w, "No available proxy", http.StatusServiceUnavailable)
		return
	}

	proxy.HandleRequest(w, r)
}

func itoa(n int) string {
	return fmt.Sprint(n)
}
