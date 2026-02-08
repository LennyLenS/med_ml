# Интеграция Python Tiler Service в composition_api

## Изменения в docker-compose.yaml

Сервис `tiler_service` теперь использует Python-версию вместо Go-версии:

```yaml
tiler_service:
  container_name: tiler_service
  build:
    context: ./tiler_python/
    dockerfile: Dockerfile
  env_file:
    - ./tiler_python/.env-docker
  ports:
    - 50080:50055
  depends_on:
    - minio
  profiles:
    - app
```

## Переменные окружения

Файл `tiler_python/.env-docker` содержит:
- `TILE_SIZE=510` - размер тайла
- `OVERLAP=1` - перекрытие тайлов
- `S3_ENDPOINT=http://minio:9000` - endpoint S3 (MinIO)
- `S3_TOKEN_ACCESS=minioadmin` - access key
- `S3_TOKEN_SECRET=minioadmin` - secret key
- `S3_BUCKET_NAME=cytology` - имя bucket
- `PORT=50055` - порт сервера
- `CACHE_DIR=/tmp/tiler_cache` - директория для кэширования

## Подключение из composition_api

В `composition-api/.env-docker` уже настроено:
```
ADAPTERS_TILERURL=http://tiler_service:50055
```

Composition API обращается к tiler_service через HTTP-клиент, который использует те же эндпоинты:
- `GET /dzi/{file_path}` - получение DZI XML
- `GET /dzi/{file_path}/files/{level}/{col}_{row}.{format}` - получение тайла

## Запуск

```bash
# Сборка и запуск всех сервисов
docker-compose --profile app up --build

# Или только tiler_service
docker-compose build tiler_service
docker-compose up tiler_service
```

## Проверка работы

```bash
# Проверка health check
curl http://localhost:50080/health

# Получение DZI XML
curl http://localhost:50080/dzi/path/to/image

# Получение тайла
curl http://localhost:50080/dzi/path/to/image/files/0/0_0.jpeg
```

## Отличия от Go-версии

1. **Логика deepZoom**: Python-версия использует правильную логику deepZoom для вычисления координат тайлов
2. **OpenSlide**: Используется openslide-python для чтения больших изображений
3. **Кэширование**: Изображения кэшируются локально после загрузки из S3
4. **API совместимость**: Полностью совместимо с Go-версией, никаких изменений в composition_api не требуется
