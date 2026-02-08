#!/usr/bin/env python3
"""
Tiler service using Python and OpenSlide
Uses deepZoom logic for tile coordinate calculation
"""
import os
import logging
import math
import tempfile
from pathlib import Path
from urllib.parse import unquote
from typing import Optional

from flask import Flask, request, Response, jsonify
import openslide
from openslide import OpenSlideError
from PIL import Image

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format='{"level":"%(levelname)s","time":"%(asctime)s","message":"%(message)s","%(name)s":"%(name)s"}',
    datefmt='%Y-%m-%d %H:%M:%S'
)
logger = logging.getLogger(__name__)

app = Flask(__name__)

# Конфигурация из переменных окружения
TILE_SIZE = int(os.getenv('TILE_SIZE', '510'))
OVERLAP = int(os.getenv('OVERLAP', '1'))
S3_ENDPOINT = os.getenv('S3_ENDPOINT', '')
S3_ACCESS_TOKEN = os.getenv('S3_TOKEN_ACCESS', '')
S3_SECRET_TOKEN = os.getenv('S3_TOKEN_SECRET', '')
S3_BUCKET_NAME = os.getenv('S3_BUCKET_NAME', 'cytology')
CACHE_DIR = os.getenv('CACHE_DIR', os.path.join(tempfile.gettempdir(), 'tiler_cache'))

# Создаем директорию кэша
os.makedirs(CACHE_DIR, exist_ok=True)

# Кэш для открытых изображений
_slide_cache = {}
_cache_lock = None

try:
    import threading
    _cache_lock = threading.Lock()
except ImportError:
    pass


def get_s3_client():
    """Получить клиент S3"""
    try:
        import boto3
        # boto3 для MinIO требует полный URL с протоколом
        # Формат: http://host:port (без пути)
        endpoint_url = S3_ENDPOINT.strip()

        # Убираем протокол, если он есть
        if endpoint_url.startswith('http://'):
            endpoint_url = endpoint_url[7:]
        elif endpoint_url.startswith('https://'):
            endpoint_url = endpoint_url[8:]

        # Убираем путь, если он есть (оставляем только host:port)
        if '/' in endpoint_url:
            endpoint_url = endpoint_url.split('/')[0]

        # Добавляем протокол http://
        endpoint_url = 'http://' + endpoint_url

        logger.info(f"Creating S3 client with endpoint: {endpoint_url}")

        # Для MinIO используем HTTP (не HTTPS)
        # boto3 автоматически определит это из протокола в endpoint_url
        s3_client = boto3.client(
            's3',
            endpoint_url=endpoint_url,
            aws_access_key_id=S3_ACCESS_TOKEN,
            aws_secret_access_key=S3_SECRET_TOKEN
        )
        return s3_client
    except ImportError:
        logger.error("boto3 not installed")
        return None
    except Exception as e:
        logger.error(f"Failed to create S3 client: {e}", exc_info=True)
        return None


def download_from_s3(s3_path: str, local_path: str) -> bool:
    """Скачать файл из S3"""
    s3_client = get_s3_client()
    if not s3_client:
        return False

    try:
        s3_client.download_file(S3_BUCKET_NAME, s3_path, local_path)
        logger.info(f"Downloaded {s3_path} to {local_path}")
        return True
    except Exception as e:
        logger.error(f"Failed to download from S3: {e}")
        return False


def get_slide_path(s3_path: str) -> Optional[str]:
    """Получить локальный путь к слайду, скачав из S3 если нужно"""
    # Создаем безопасное имя файла из пути
    safe_name = s3_path.replace('/', '_').replace('\\', '_')
    local_path = os.path.join(CACHE_DIR, safe_name)

    # Если файл уже есть, возвращаем путь
    if os.path.exists(local_path):
        return local_path

    # Скачиваем из S3
    if download_from_s3(s3_path, local_path):
        return local_path

    return None


def get_slide(s3_path: str) -> Optional[openslide.OpenSlide]:
    """Получить OpenSlide объект для изображения"""
    # Проверяем кэш
    if _cache_lock:
        with _cache_lock:
            if s3_path in _slide_cache:
                return _slide_cache[s3_path]

    # Получаем локальный путь
    local_path = get_slide_path(s3_path)
    if not local_path:
        return None

    try:
        slide = openslide.OpenSlide(local_path)

        # Сохраняем в кэш
        if _cache_lock:
            with _cache_lock:
                _slide_cache[s3_path] = slide

        return slide
    except OpenSlideError as e:
        logger.error(f"Failed to open slide {s3_path}: {e}")
        return None
    except Exception as e:
        logger.error(f"Unexpected error opening slide {s3_path}: {e}")
        return None


def get_dzi_xml(slide: openslide.OpenSlide) -> str:
    """Генерировать DZI XML для изображения"""
    width, height = slide.dimensions

    # Вычисляем количество уровней как в OpenSeadragon
    # maxLevel = floor(log2(maxDim)) - это гарантирует, что на level 0 будет >= 1px
    max_dim = max(width, height)
    if max_dim <= 0:
        max_level = 0
    else:
        max_level = int(math.floor(math.log2(max_dim)))

    xml = f'''<?xml version="1.0" encoding="UTF-8"?>
<Image xmlns="http://schemas.microsoft.com/deepzoom/2008"
       TileSize="{TILE_SIZE}"
       Overlap="{OVERLAP}"
       Format="jpeg"
       ServerFormat="Default">
  <Size Width="{width}" Height="{height}"/>
</Image>'''
    return xml


@app.route('/dzi/<path:file_path>', methods=['GET'])
def get_dzi(file_path: str):
    """Получить DZI XML для изображения"""
    try:
        decoded_path = unquote(file_path)
        logger.info(f"GetDZI request: {decoded_path}")

        slide = get_slide(decoded_path)
        if not slide:
            return jsonify({"error": "Failed to open image"}), 500

        dzi_xml = get_dzi_xml(slide)

        return Response(
            dzi_xml,
            mimetype='application/xml',
            status=200
        )
    except Exception as e:
        logger.error(f"GetDZI error: {e}", exc_info=True)
        return jsonify({"error": str(e)}), 500


@app.route('/dzi/<path:file_path>/files/<int:level>/<tile_coords>.<format>', methods=['GET'])
def get_tile(file_path: str, level: int, tile_coords: str, format: str):
    """Получить тайл изображения"""
    try:
        decoded_path = unquote(file_path)

        # Парсим координаты тайла (формат: col_row)
        try:
            col, row = map(int, tile_coords.split('_'))
        except ValueError:
            return jsonify({"error": "Invalid tile coordinates"}), 400

        logger.info(f"GetTile request: {decoded_path}, level={level}, col={col}, row={row}, format={format}")

        slide = get_slide(decoded_path)
        if not slide:
            return jsonify({"error": "Failed to open image"}), 500

        # Вычисляем количество уровней (как в deepZoom/OpenSeadragon)
        width, height = slide.dimensions
        max_dim = max(width, height)
        if max_dim <= 0:
            max_level = 0
        else:
            max_level = int(math.floor(math.log2(max_dim)))

        # Вычисляем координаты тайла на уровне DZI (как в deepZoom)
        tile_x = col * TILE_SIZE
        tile_y = row * TILE_SIZE

        # Вычисляем масштаб уровня (как в deepZoom)
        # В OpenSeadragon: maxLevel = полное разрешение, level 0 = минимальный
        # scale = 2^(maxLevel - level) - это масштаб уровня относительно полного разрешения
        dzi_scale = 2.0 ** (max_level - level)

        # Координаты в полном разрешении (как в deepZoom)
        source_x = int((tile_x - OVERLAP) * dzi_scale)
        source_y = int((tile_y - OVERLAP) * dzi_scale)
        source_width = int((TILE_SIZE + 2 * OVERLAP) * dzi_scale)
        source_height = int((TILE_SIZE + 2 * OVERLAP) * dzi_scale)

        # Ограничиваем границами изображения
        source_x = max(0, source_x)
        source_y = max(0, source_y)
        source_width = min(source_width, width - source_x)
        source_height = min(source_height, height - source_y)

        if source_width <= 0 or source_height <= 0:
            return jsonify({"error": "Tile coordinates out of bounds"}), 404

        # Находим подходящий уровень OpenSlide
        level_count = slide.level_count
        best_level = 0
        best_scale_diff = float('inf')

        # Желаемый масштаб для уровня DZI
        desired_scale = 1.0 / dzi_scale  # Масштаб уровня OpenSlide относительно полного разрешения

        for i in range(level_count):
            level_dimensions = slide.level_dimensions[i]
            level_scale = width / level_dimensions[0]
            scale_diff = abs(level_scale - desired_scale)
            if scale_diff < best_scale_diff:
                best_scale_diff = scale_diff
                best_level = i

        # Координаты на уровне OpenSlide
        level_dimensions = slide.level_dimensions[best_level]
        os_level_scale = width / level_dimensions[0]

        level_x = int(source_x / os_level_scale)
        level_y = int(source_y / os_level_scale)
        level_w = int(source_width / os_level_scale)
        level_h = int(source_height / os_level_scale)

        # Ограничиваем размер области для чтения
        max_read_size = TILE_SIZE * 10
        if level_w > max_read_size:
            level_w = max_read_size
        if level_h > max_read_size:
            level_h = max_read_size

        # Читаем область
        tile_image = slide.read_region((level_x, level_y), best_level, (level_w, level_h))

        # Конвертируем RGBA в RGB
        if tile_image.mode == 'RGBA':
            # Создаем белый фон
            rgb_image = Image.new('RGB', tile_image.size, (255, 255, 255))
            rgb_image.paste(tile_image, mask=tile_image.split()[3])  # Используем альфа-канал как маску
            tile_image = rgb_image

        # Масштабируем до нужного размера тайла
        target_size = (TILE_SIZE + 2 * OVERLAP, TILE_SIZE + 2 * OVERLAP)
        if tile_image.size != target_size:
            tile_image = tile_image.resize(target_size, Image.Resampling.LANCZOS)

        # Кодируем в нужный формат
        import io
        output = io.BytesIO()

        if format.lower() in ('jpeg', 'jpg'):
            tile_image.save(output, format='JPEG', quality=85)
            mimetype = 'image/jpeg'
        elif format.lower() == 'png':
            tile_image.save(output, format='PNG')
            mimetype = 'image/png'
        else:
            return jsonify({"error": f"Unsupported format: {format}"}), 400

        output.seek(0)

        return Response(
            output.getvalue(),
            mimetype=mimetype,
            status=200
        )

    except Exception as e:
        logger.error(f"GetTile error: {e}", exc_info=True)
        return jsonify({"error": str(e)}), 500


@app.route('/health', methods=['GET'])
def health():
    """Health check endpoint"""
    return jsonify({"status": "ok"}), 200


if __name__ == '__main__':
    port = int(os.getenv('PORT', '50055'))
    app.run(host='0.0.0.0', port=port, debug=False)
