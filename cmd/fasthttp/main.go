package main

import (
	"log"

	"github.com/jwhittle933/gateway"
	"github.com/jwhittle933/gateway/proxy"

	"github.com/valyala/fasthttp"
)

func main() {
	g := gateway.New(
		gateway.NewConfig("/yahoo", proxy.New("https://yahoo.com", nil)),
		gateway.NewConfig("/bing", proxy.New("https://bing.com", nil)),
	)

	g.Proxy(
		"/google",
		proxy.New(
			"https://google.com",
			nil,
			proxy.NewConfig(proxy.Default, fasthttp.MethodGet, "", nil),
		),
	)

	s := &fasthttp.Server{Handler: g.Handle}
	if err := s.ListenAndServe("8080"); err != nil {
		log.Fatalf("Error closing server: %s", err.Error())
	}
}
