# Tiler Service (Python)

Сервис для генерации тайлов изображений с использованием Python и OpenSlide.
Использует логику deepZoom для вычисления координат тайлов.

## Зависимости

- Python 3.11+
- OpenSlide (libopenslide0, libopenslide-dev)
- Flask
- openslide-python
- Pillow
- boto3 (для S3)

## Переменные окружения

- `TILE_SIZE` - размер тайла (по умолчанию: 510)
- `OVERLAP` - перекрытие тайлов (по умолчанию: 1)
- `S3_ENDPOINT` - endpoint S3 хранилища
- `S3_TOKEN_ACCESS` - access token для S3
- `S3_TOKEN_SECRET` - secret token для S3
- `S3_BUCKET_NAME` - имя bucket в S3 (по умолчанию: cytology)
- `CACHE_DIR` - директория для кэширования изображений (по умолчанию: /tmp/tiler_cache)
- `PORT` - порт сервера (по умолчанию: 50055)

## API

### GET /dzi/{file_path}

Возвращает DZI XML для изображения.

### GET /dzi/{file_path}/files/{level}/{col}_{row}.{format}

Возвращает тайл изображения.

Параметры:
- `file_path` - путь к файлу в S3
- `level` - уровень масштаба (0 = минимальный, maxLevel = полное разрешение)
- `col` - номер колонки тайла
- `row` - номер строки тайла
- `format` - формат изображения (jpeg, png)

### GET /health

Health check endpoint.

## Запуск

```bash
# Установка зависимостей
pip install -r requirements.txt

# Запуск сервера
python main.py
```

Или через Docker:

```bash
docker build -t tiler-python .
docker run -p 50055:50055 \
  -e S3_ENDPOINT=... \
  -e S3_TOKEN_ACCESS=... \
  -e S3_TOKEN_SECRET=... \
  tiler-python
```
