package proxy

import (
	"fmt"
	"net/http"
	"path"

	"github.com/jwhittle933/gonads/result"
	"github.com/valyala/fasthttp"
)

const (
	Default string = "default"
)

// Connection represents a connection from one service to another
type Connection struct {
	BaseUrl        string
	DefaultHeaders *fasthttp.RequestHeader
	factories      RequestFactories
	client         *fasthttp.Client
}

type RequestFactories map[string]FactoryFunc
type FactoryFunc func([]byte, ...interface{}) *fasthttp.Request

type Config struct {
	Name          string
	Method        string
	FormattedPath string
	Headers       map[string]string
}

func New(baseURL string, defaultHeaders map[string]string, options ...Config) *Connection {
	return (&Connection{BaseUrl: baseURL, client: &fasthttp.Client{}, factories: RequestFactories{}}).
		withFactories(options...).
		withDefaultHeaders(defaultHeaders)
}

func NewConfig(name, method, formattedPath string, headers map[string]string) Config {
	return Config{name, method, formattedPath, headers}
}

// Submit performs the named request supplied to the factory
// and returns the result.Result
//
// Params: If the request needs path params beyond the base url supplied at creation,
// the base url must be a valid format string, i.e., `/resource/%d`, and the integer
// path param supplied when Submit is called
//
// `name` is the request name given when created.
// `body` is the byte-encoded request body, `nil` if no body is required
// `with` is a variadic, ordered sequence of path params that should correspond
// to the format string used. These can be split ("first", "second"), or
// combined ("first/second").
func (c Connection) Submit(name string, body []byte, with ...interface{}) result.Result {
	req := c.factories.Request(name, body, with...)
	res := fasthttp.AcquireResponse()

	err := c.client.Do(req, res)
	return result.Handle(res, err)
}

func (c Connection) SubmitDefault(body []byte) result.Result {
	return c.Submit("", body)
}

func (c *Connection) withFactories(configs ...Config) *Connection {
	if len(configs) == 0 {
		c.factories.Set("default", c.newFactoryFunc(Config{"default", http.MethodGet, "", nil}))
		return c
	}

	for _, conf := range configs {
		c.factories.Set(conf.Name, c.newFactoryFunc(conf))
	}

	return c
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

func (c *Connection) newFactoryFunc(config Config) FactoryFunc {
	return func(body []byte, pathParams ...interface{}) *fasthttp.Request {
		req := fasthttp.AcquireRequest()
		c.DefaultHeaders.CopyTo(&req.Header)

		req.Header.SetMethod(config.Method)
		req.SetRequestURI(path.Join(c.BaseUrl, applyPathOpts(config.FormattedPath, pathParams...)))
		req.SetBody(body)

		if config.Headers != nil {
			for k, v := range config.Headers {
				req.Header.Set(k, v)
			}
		}

		return req
	}
}

func (rf RequestFactories) Set(id string, ff FactoryFunc) {
	rf[id] = ff
}

func (rf RequestFactories) Request(name string, body []byte, with ...interface{}) *fasthttp.Request {
	if name == "" {
		return rf["default"](body)
	}

	if dispatcher, ok := rf[name]; ok {
		return dispatcher(body, with...)
	}

	return rf["default"](body)
}

func applyPathOpts(basePath string, params ...interface{}) string {
	if len(params) == 0 {
		return basePath
	}

	return fmt.Sprintf(basePath, params...)
}

// Proxy result pipeline result.Binder helper
func Proxy(ctx *fasthttp.RequestCtx) result.Binder {
	return func(val interface{}) result.Result {
		res := val.(*fasthttp.Response)
		res.CopyTo(&ctx.Response)

		return result.Wrap(res)
	}
}
