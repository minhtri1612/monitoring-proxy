FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY main.go .
RUN CGO_ENABLED=0 go build -o demo-api main.go

FROM alpine:3.20
COPY --from=builder /app/demo-api /demo-api
EXPOSE 8080
ENTRYPOINT ["/demo-api"]
