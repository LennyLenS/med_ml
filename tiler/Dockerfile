FROM golang:alpine

# Устанавливаем git, необходимый для загрузки зависимостей из git репозиториев
RUN apk add --no-cache git

WORKDIR /service

# Копируем go.mod и go.sum для кэширования слоев Docker
COPY go.mod go.sum ./

# Загружаем существующие зависимости (если go.sum есть)
RUN go mod download || true

# Копируем весь исходный код
COPY . .

# Получаем недостающие зависимости и обновляем go.mod/go.sum
RUN go mod tidy

# Собираем приложение
RUN go build -o bin/tiler cmd/service/main.go

EXPOSE 50055

CMD ["./bin/tiler"]
