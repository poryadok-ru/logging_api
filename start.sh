#!/bin/bash

# Остановить и удалить старый контейнер, если он существует
docker stop logging_api 2>/dev/null || true
docker rm logging_api 2>/dev/null || true

# Собрать образ
echo "Сборка Docker образа..."
docker build -t logging_api:latest .

# Запустить контейнер
echo "Запуск контейнера..."
docker run -d \
  --name logging_api \
  --restart unless-stopped \
  --network host \
  --env-file .env \
  logging_api:latest

echo "Контейнер запущен!"
echo "Проверьте логи: docker logs -f logging_api"
echo "API доступен на http://localhost:8080"
echo "Swagger: http://localhost:8080/swagger/index.html"

