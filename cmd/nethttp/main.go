package main

import (
	"github.com/jwhittle933/gateway"
	"github.com/jwhittle933/gateway/proxy"
	"log"
	"net/http"
)

func main() {
	g := gateway.New()

	g.Proxy(
		"/things",
		proxy.New(
			"https://my-domain.com",
			map[string]string{"accept": "application/json"},
			proxy.NewConfig(proxy.Default, http.MethodGet, "/things", nil),
		),
	)

	if err := http.ListenAndServe("8080", g); err != nil {
		log.Fatalf("Error closing server: %s", err.Error())
	}
}
