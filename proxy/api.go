package proxy

import (
	"fmt"
	"github.com/jwhittle933/gateway/proxy/result"
	"path"

	"github.com/valyala/fasthttp"
)

// ApiResult local alias for result.Result
type ApiResult result.Result

type Connector interface {
	Submit(name string, body []byte, with ...interface{}) ApiResult
	SubmitDefault(body []byte) ApiResult
}

type Connection struct {
	BaseUrl        string
	DefaultHeaders *fasthttp.RequestHeader
	factories      RequestFactories
	client         *fasthttp.Client
}

type RequestFactories map[string]FactoryFunc
type FactoryFunc func([]byte, ...interface{}) *fasthttp.Request

type Option struct {
	Name          string
	Method        string
	FormattedPath string
	Headers       map[string]string
}

func New(client *fasthttp.Client, baseURL string, defaultHeaders map[string]string, options ...Option) *Connection {
	return (&Connection{BaseUrl: baseURL, client: client, factories: RequestFactories{}}).
		withFactories(options...).
		withDefaultHeaders(defaultHeaders)
}

func (c *Connection) withFactories(options ...Option) *Connection {

}

func (c *Connection) withDefaultHeaders(headers map[string]string) *Connection {
	if headers == nil {
		return c
	}

	c.DefaultHeaders = &fasthttp.RequestHeader{}
	for k, v := range headers {
		c.DefaultHeaders.Set(k, v)
	}

	return c
}

func (c *Connection) newFactoryFunc(option Option) FactoryFunc {
	return func(body []byte, pathParams ...interface{}) *fasthttp.Request {
		req := fasthttp.AcquireRequest()
		c.DefaultHeaders.CopyTo(&req.Header)

		req.Header.SetMethod(option.Method)
		req.SetRequestURI(path.Join(c.BaseUrl, applyPathOpts(option.FormattedPath, pathParams...)))
	}
}

func applyPathOpts(basePath string, params ...interface{}) string {
	if len(params) == 0 {
		return basePath
	}

	return fmt.Sprintf(basePath, params...)
}
