#!/usr/bin/env python3
"""
Tiler service using Python, OpenSlide and DeepZoomGenerator
Simple version with disk cache for files and in-memory cache for generators
"""
import os
import logging
import tempfile
import threading
from io import BytesIO
from urllib.parse import unquote
from typing import Optional

from flask import Flask, Response, jsonify
from flask_cors import CORS
from openslide import OpenSlide
from openslide.deepzoom import DeepZoomGenerator
from cachetools import TTLCache

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format='{"level":"%(levelname)s","time":"%(asctime)s","message":"%(message)s"}',
    datefmt='%Y-%m-%d %H:%M:%S'
)
logger = logging.getLogger(__name__)

app = Flask(__name__)

# Настройка CORS
CORS(app, resources={r"/*": {"origins": "*"}})

# Конфигурация из переменных окружения
TILE_SIZE = int(os.getenv('TILE_SIZE', '510'))
OVERLAP = int(os.getenv('OVERLAP', '1'))
S3_ENDPOINT = os.getenv('S3_ENDPOINT', '')
S3_ACCESS_TOKEN = os.getenv('S3_TOKEN_ACCESS', '')
S3_SECRET_TOKEN = os.getenv('S3_TOKEN_SECRET', '')
S3_BUCKET_NAME = os.getenv('S3_BUCKET_NAME', 'cytology')
CACHE_DIR = os.getenv('CACHE_DIR', os.path.join(tempfile.gettempdir(), 'tiler_cache'))
PORT = int(os.getenv('PORT', '50055'))

# Создаем директорию кэша
os.makedirs(CACHE_DIR, exist_ok=True)
logger.info(f"Cache directory: {CACHE_DIR}")

# Кэш для DeepZoomGenerator (в памяти)
_generator_cache = TTLCache(maxsize=200, ttl=3600)
_cache_lock = threading.Lock()


def get_s3_client():
    """Получить клиент S3"""
    try:
        import boto3

        # Обработка endpoint URL
        endpoint_url = S3_ENDPOINT.strip()

        # Убираем протокол, если есть
        if endpoint_url.startswith('http://'):
            endpoint_url = endpoint_url[7:]
        elif endpoint_url.startswith('https://'):
            endpoint_url = endpoint_url[8:]

        # Убираем путь, если есть
        if '/' in endpoint_url:
            endpoint_url = endpoint_url.split('/')[0]

        # Добавляем протокол http://
        endpoint_url = 'http://' + endpoint_url

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
    """Скачать файл из S3 на диск"""
    s3_client = get_s3_client()
    if not s3_client:
        return False

    try:
        logger.info(f"Downloading {s3_path} from bucket {S3_BUCKET_NAME}")
        s3_client.download_file(S3_BUCKET_NAME, s3_path, local_path)
        logger.info(f"Downloaded to {local_path}")
        return True
    except Exception as e:
        logger.error(f"Failed to download from S3: {e}", exc_info=True)
        return False


def get_slide_path(s3_path: str) -> Optional[str]:
    """Получить локальный путь к слайду, скачав из S3 если нужно (кэш на диске)"""
    # Создаем безопасное имя файла
    safe_name = s3_path.replace('/', '_').replace('\\', '_')
    local_path = os.path.join(CACHE_DIR, safe_name)

    # Если файл уже есть на диске, возвращаем путь
    if os.path.exists(local_path):
        return local_path

    # Скачиваем из S3
    if download_from_s3(s3_path, local_path):
        return local_path

    return None


def get_slide_generator(s3_path: str) -> Optional[DeepZoomGenerator]:
    """Получить DeepZoomGenerator для изображения (кэш в памяти)"""
    # Проверяем кэш генераторов
    with _cache_lock:
        generator = _generator_cache.get(s3_path)
        if generator is not None:
            return generator

    # Получаем локальный путь (кэш на диске)
    local_path = get_slide_path(s3_path)
    if not local_path:
        logger.error(f"Failed to get local path for {s3_path}")
        return None

    try:
        # Открываем слайд
        slide = OpenSlide(local_path)

        # Создаем генератор DeepZoom
        generator = DeepZoomGenerator(
            slide,
            tile_size=TILE_SIZE,
            overlap=OVERLAP,
            limit_bounds=True
        )

        # Сохраняем в кэш генераторов
        with _cache_lock:
            _generator_cache[s3_path] = generator

        logger.info(f"Created DeepZoomGenerator for {s3_path}")
        return generator
    except Exception as e:
        logger.error(f"Failed to create DeepZoomGenerator for {s3_path}: {e}", exc_info=True)
        return None


@app.route('/dzi/<path:file_path>', methods=['GET'])
def get_dzi(file_path: str):
    """Получить DZI XML для изображения"""
    try:
        decoded_path = unquote(file_path)
        logger.info(f"GetDZI request: {decoded_path}")

        generator = get_slide_generator(decoded_path)
        if not generator:
            return jsonify({"error": "Failed to open image"}), 500

        # DeepZoomGenerator сам генерирует правильный DZI XML
        dzi_xml = generator.get_dzi("jpeg")

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

        generator = get_slide_generator(decoded_path)
        if not generator:
            return jsonify({"error": "Failed to open image"}), 500

        # DeepZoomGenerator сам правильно вычисляет координаты и масштаб
        try:
            tile = generator.get_tile(level, (col, row))
        except ValueError as e:
            # Тайл не существует (out of bounds)
            logger.warn(f"Tile not found: level={level}, col={col}, row={row}, error={e}")
            return jsonify({"error": "Tile not found"}), 404

        # Сохраняем в буфер
        buf = BytesIO()
        tile.save(buf, format.lower())

        # Определяем MIME type
        mimetype_map = {
            'jpeg': 'image/jpeg',
            'jpg': 'image/jpeg',
            'png': 'image/png'
        }
        mimetype = mimetype_map.get(format.lower(), 'image/jpeg')

        return Response(
            buf.getvalue(),
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
    logger.info(f"Starting tiler service on port {PORT}")
    logger.info(f"Tile size: {TILE_SIZE}, Overlap: {OVERLAP}")
    app.run(host='0.0.0.0', port=PORT, debug=False)
