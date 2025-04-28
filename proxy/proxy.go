package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type HttpProxy struct {
	Adress      url.URL
	connections int
	rmx         sync.RWMutex
	isAlive     bool
	proxy       httputil.ReverseProxy
}

func (hp *HttpProxy) HandleRequest(w http.ResponseWriter, r *http.Request) {
	hp.proxy.ServeHTTP(w, r)
}

func (hp *HttpProxy) GetConnection() int {
	return hp.connections
}

func (hp *HttpProxy) GetAdress() string {
	return hp.Adress.RequestURI()
}

func (hp *HttpProxy) SetAlive(isAlive bool) {
	hp.isAlive = isAlive
}

func (hp *HttpProxy) GetStatusIsAlive() bool {
	return hp.isAlive
}

func (hp *HttpProxy) IsAlive() bool {
	hp.rmx.Lock()
	defer hp.rmx.Unlock()

	resStatus, err := hp.ping()
	if err != nil {
		return false
	}

	return resStatus >= 200 && resStatus <= 400
}

func (hp *HttpProxy) ping() (int, error) {
	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest("HEAD", hp.GetAdress(), nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()

	return resp.StatusCode, nil
}
