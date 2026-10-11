package main

import (
	"io"
	"log"
	"net/http"
)

const page = `<!doctype html><html><head><meta charset="utf-8"><title>Prometheus Proxy Demo</title><style>body{font:16px system-ui;margin:40px;background:#f5f7fb;color:#162033}.grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:20px}.card{padding:22px;border-radius:12px;background:#fff;box-shadow:0 2px 10px #dde2ee}pre{white-space:pre-wrap;background:#101828;color:#c8f7d4;padding:14px;border-radius:8px}button{padding:10px 14px}</style></head><body><h1>Prometheus Proxy: tenant isolation</h1><p>Each card asks the same metric through a different API key. Refresh after running k6.</p><button onclick="load()">Refresh</button><div class="grid"><section class="card"><h2>Production</h2><pre id="production">Loading…</pre></section><section class="card"><h2>Staging</h2><pre id="staging">Loading…</pre></section></div><script>async function load(){for(const e of ['production','staging']){const r=await fetch('/api/'+e);document.getElementById(e).textContent=JSON.stringify(await r.json(),null,2)}}load()</script></body></html>`

func proxyQuery(key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest(http.MethodGet, "http://proxy:8000/api/v1/query?query=demo_http_requests_total", nil)
		if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
		req.Header.Set("Authorization", "Basic demo-proxy-auth")
		req.Header.Set("PROMETHEUS-API-KEY", key)
		resp, err := http.DefaultClient.Do(req)
		if err != nil { http.Error(w, err.Error(), http.StatusBadGateway); return }
		defer resp.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "text/html; charset=utf-8"); io.WriteString(w, page) })
	http.HandleFunc("/api/production", proxyQuery("demo-production-key"))
	http.HandleFunc("/api/staging", proxyQuery("demo-staging-key"))
	log.Fatal(http.ListenAndServe(":3001", nil))
}
