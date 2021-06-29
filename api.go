package gateway

import (
	"github.com/jwhittle933/gonads/result"
	"github.com/valyala/fasthttp"
	"net/http"
)

type Connector interface {
	Submit(name string, body []byte, with ...interface{}) result.Result
	SubmitDefault(body []byte) result.Result
}

type Shuttle interface {
	Go()
}

// Gateway represents an 1-to-1 reverse proxy
// to other services
type Gateway struct {
	connections map[Path]Connector
}

type Config struct {
	Path string
	Connector Connector
}

// New factory func creates a new Gateway
func New(configs ...Config) *Gateway {
	conns := make(map[Path]Connector)
	for _, config := range configs {
		conns[NewPath(config.Path)] = config.Connector
	}

	return &Gateway{connections: conns}
}

func NewConfig(path string, c Connector) Config {
	return Config{path, c}
}

func (g *Gateway) Proxy(path string, c Connector) {
	g.connections[NewPath(path)] = c
}

// ServeHTTP satisfies the http.Handler interface
func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// todo
}

// Handle for use with fasthttp
func (g *Gateway) Handle(ctx *fasthttp.RequestCtx) {
	//
}
