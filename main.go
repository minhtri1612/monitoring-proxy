package main

import (
	"log"
	"net/http"

	"github.com/galaxyed/prometheus-proxy/internal/conf"
	"github.com/galaxyed/prometheus-proxy/internal/server"
)

func main() {
	// Read Config
	cfgPath, err := server.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := conf.NewConfig(cfgPath)

	// Initialize the reverse proxy with the Prometheus URL from the configuration.
	proxy, err := server.NewProxy(cfg.PrometheusURL, cfg)
	if err != nil {
		panic(err)
	}
	log.Println(cfg.Policies[0].Name)

	// handle all requests to your server using the proxy
	http.HandleFunc("/", server.ProxyRequestHandler(proxy, cfg))
	log.Println("Server Started")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}
