package proxy

import (
	"fmt"
	"net/http"
	"path"

	"github.com/valyala/fasthttp"

	"github.com/jwhittle933/gateway/proxy/result"
)

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
	if err := c.client.Do(req, res); err != nil {
		return result.WrapErr(err)
	}

	return result.Wrap(res)
}

func (c *Connection) withFactories(options ...Option) *Connection {
	if len(options) == 0 {
		c.factories.Set("default", c.newFactoryFunc(Option{"default", http.MethodGet, "", nil}))
		return c
	}

	for _, op := range options {
		c.factories.Set(op.Name, c.newFactoryFunc(op))
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

func (c *Connection) newFactoryFunc(option Option) FactoryFunc {
	return func(body []byte, pathParams ...interface{}) *fasthttp.Request {
		req := fasthttp.AcquireRequest()
		c.DefaultHeaders.CopyTo(&req.Header)

		req.Header.SetMethod(option.Method)
		req.SetRequestURI(path.Join(c.BaseUrl, applyPathOpts(option.FormattedPath, pathParams...)))
		req.SetBody(body)

		if option.Headers != nil {
			for k, v := range option.Headers {
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

func Proxy(ctx *fasthttp.RequestCtx) result.Binder {
	return func(val interface{}) result.Result {
		res := val.(*fasthttp.Response)
		res.CopyTo(&ctx.Response)

		return result.Wrap(res)
	}
}
