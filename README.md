# prometheus-proxy

Reverse proxy for Prometheus that adds label filters based on a caller's
`PROMETHEUS-API-KEY`. It lets a dashboard or application query only the
metrics assigned to its policy instead of accessing Prometheus directly.

## Run the local demo

The included Compose lab starts Prometheus, two node exporters, this proxy,
a small dashboard, and two demo APIs. Prometheus labels the targets as
`cluster=production` and `cluster=staging`.

```sh
docker compose up --build
```

Open:

- Grafana: http://localhost:3000 (login `admin` / `admin`)
- Demo dashboard: http://localhost:3001
- Proxy: http://localhost:8000
- Prometheus (for comparison only): http://localhost:9090

In Grafana Explore, select **Prometheus via Proxy (production)** and run:

```promql
up
```

The proxy adds `cluster="production"`, so only the production target appears.
To test the proxy without Grafana:

```sh
curl -G http://localhost:8000/api/v1/query \
  -H 'Authorization: Basic demo-proxy-auth' \
  -H 'PROMETHEUS-API-KEY: demo-production-key' \
  --data-urlencode 'query=up'
```

Stop the lab with `docker compose down`.

## Load-test tenant isolation

The demo APIs expose `GET`, `POST`, and `PUT /api`. Run the supplied k6 test
with 1,000 virtual users for 30 seconds:

```sh
docker compose --profile load run --rm k6
```

After it completes, refresh http://localhost:3001. The production card must
only contain time series labeled `cluster="production"`; the staging card must
only contain `cluster="staging"`. Requests without valid proxy headers now
receive `403 Forbidden` rather than being forwarded to Prometheus.

## Configuration

`conf.yml` defines the upstream `PrometheusURL`, the internal authorization
header, and the label policies. Do not commit real credentials or API keys in
this file for production.
