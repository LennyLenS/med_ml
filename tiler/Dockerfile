FROM golang:alpine

# Устанавливаем зависимости для libvips и OpenSlide
# vips-dev содержит заголовочные файлы для компиляции
# vips содержит runtime библиотеку
# openslide-dev содержит заголовочные файлы для OpenSlide
# openslide содержит runtime библиотеку OpenSlide
# gcc, musl-dev нужны для CGO компиляции
# pkg-config нужен для поиска библиотек через pkg-config
RUN apk add --no-cache \
	git \
	gcc \
	musl-dev \
	pkgconfig \
	vips-dev \
	vips \
	openslide-dev \
	openslide

WORKDIR /service

# Копируем go.mod и go.sum для кэширования слоев Docker
COPY go.mod go.sum ./

# Загружаем существующие зависимости (если go.sum есть)
RUN go mod download || true

# Копируем весь исходный код
COPY . .

# Получаем недостающие зависимости и обновляем go.mod/go.sum
RUN go mod tidy

# Собираем приложение с CGO (требуется для govips/libvips)
ENV CGO_ENABLED=1
RUN go build -o bin/tiler cmd/service/main.go

EXPOSE 50055

CMD ["./bin/tiler"]
