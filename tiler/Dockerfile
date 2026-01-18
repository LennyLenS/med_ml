FROM golang:alpine

WORKDIR /service

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем и обновляем зависимости
RUN go mod download && go mod tidy

# Копируем остальные файлы
COPY . .

# Собираем приложение
RUN go build -o bin/tiler cmd/service/main.go

EXPOSE 50055

CMD ["./bin/tiler"]
