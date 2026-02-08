FROM python:3.11-slim

# Устанавливаем системные зависимости для OpenSlide
RUN apt-get update && apt-get install -y --no-install-recommends \
    libopenslide0 \
    libopenslide-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /service

# Копируем requirements.txt
COPY requirements.txt .

# Устанавливаем Python зависимости
RUN pip install --no-cache-dir -r requirements.txt

# Копируем исходный код
COPY . .

# Открываем порт
EXPOSE 50055

# Запускаем приложение
CMD ["python", "main.py"]
