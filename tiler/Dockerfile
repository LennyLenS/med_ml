FROM golang:alpine

# Устанавливаем git, необходимый для загрузки зависимостей из git репозиториев
RUN apk add --no-cache git

WORKDIR /service

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Получаем последнюю версию библиотеки и загружаем зависимости
RUN go get github.com/sunshineplan/tiff@latest
RUN go mod download
RUN go mod tidy

# Копируем остальные файлы
COPY . .

# Собираем приложение
RUN go build -o bin/tiler cmd/service/main.go

EXPOSE 50055

CMD ["./bin/tiler"]
