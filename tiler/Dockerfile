# Используем стабильный тег golang:bookworm или golang:1.23-bookworm
# Если проблемы с сетью, попробуйте: golang:latest или golang:1.23
FROM golang:bookworm

# Устанавливаем зависимости для libvips и OpenSlide
# libvips-dev содержит заголовочные файлы для компиляции
# libvips-tools содержит утилиты
# libopenslide0 содержит runtime библиотеку OpenSlide
# libopenslide-dev содержит заголовочные файлы для OpenSlide
# gcc, g++, pkg-config нужны для CGO компиляции
RUN apt-get update && apt-get install -y --no-install-recommends \
	git \
	gcc \
	g++ \
	pkg-config \
	libvips-dev \
	libvips-tools \
	libopenslide0 \
	libopenslide-dev \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /service

# Копируем go.mod и go.sum для кэширования слоев Docker
COPY go.mod go.sum ./

# Загружаем существующие зависимости (если go.sum есть)
RUN go mod download || true

# Копируем весь исходный код
COPY . .

# Получаем недостающие зависимости и обновляем go.mod/go.sum
RUN go mod tidy

# Собираем приложение с CGO (требуется для OpenSlide и libvips)
ENV CGO_ENABLED=1
RUN go build -o bin/tiler cmd/service/main.go

EXPOSE 50055

CMD ["./bin/tiler"]
