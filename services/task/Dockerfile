FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o task ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/task .
EXPOSE 8082
ENTRYPOINT ["./task"]
