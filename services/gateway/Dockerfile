FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o gateway ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/gateway .
EXPOSE 8080
ENTRYPOINT ["./gateway"]
