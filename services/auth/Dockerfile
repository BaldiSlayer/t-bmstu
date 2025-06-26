FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o auth ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/auth .
EXPOSE 8082
ENTRYPOINT ["./auth"]
