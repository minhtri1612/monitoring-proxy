package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

var requests [3]uint64

func metricIndex(method string) int {
	switch method {
	case http.MethodGet:
		return 0
	case http.MethodPost:
		return 1
	default:
		return 2
	}
}

func main() {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "unknown"
	}

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost && r.Method != http.MethodPut {
			http.Error(w, "only GET, POST, and PUT are supported", http.StatusMethodNotAllowed)
			return
		}
		atomic.AddUint64(&requests[metricIndex(r.Method)], 1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"environment":%q,"method":%q,"ok":true}`+"\n", environment, r.Method)
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fmt.Fprintln(w, "# HELP demo_http_requests_total Requests received by the load-test demo API.")
		fmt.Fprintln(w, "# TYPE demo_http_requests_total counter")
		for index, method := range []string{"GET", "POST", "PUT"} {
			fmt.Fprintf(w, "demo_http_requests_total{method=%q,service=%q} %d\n", method, "demo-api", atomic.LoadUint64(&requests[index]))
		}
	})

	log.Printf("demo API for %s listening on :8080", environment)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
