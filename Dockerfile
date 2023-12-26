FROM golang:latest

WORKDIR /app

# Копируем все файлы из текущего каталога (где находится Dockerfile) внутрь контейнера
COPY . .

RUN go mod tidy

# Сборка приложения
RUN go build -o main ./cmd/

EXPOSE 8080

# Запуск приложения
CMD ["./main"]
