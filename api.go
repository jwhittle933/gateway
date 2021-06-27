package gatway

import (
	"github.com/jwhittle933/gateway/proxy/result"
)

type Path string

// ApiResult local alias for result.Result
type ApiResult result.Result

type Connector interface {
	Submit(name string, body []byte, with ...interface{}) ApiResult
	SubmitDefault(body []byte) ApiResult
}

// Gateway represents an 1-to-1 reverse proxy
// to other services
type Gateway struct {
	Connections map[Path]Connector
}

func (g *Gateway) Proxy() {
	//
}

func (g *Gateway) Intercept() {
	//
}