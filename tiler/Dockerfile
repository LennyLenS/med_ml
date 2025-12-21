FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go mod файлы
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tiler ./cmd/service

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/tiler .

EXPOSE 50055

CMD ["./tiler"]
