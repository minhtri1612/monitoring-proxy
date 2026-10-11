FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY main.go .
RUN CGO_ENABLED=0 go build -o dashboard main.go

FROM alpine:3.20
COPY --from=builder /app/dashboard /dashboard
EXPOSE 3001
ENTRYPOINT ["/dashboard"]
