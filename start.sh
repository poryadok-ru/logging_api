#!/bin/bash

# Проверка наличия .env файла
if [ ! -f .env ]; then
    echo "⚠️  Файл .env не найден!"
    echo "Скопируйте .env.example в .env и настройте параметры"
    exit 1
fi

# Остановить и удалить старый контейнер, если существует
docker stop logging_api 2>/dev/null || true
docker rm logging_api 2>/dev/null || true

# Запустить контейнер
echo "Запуск контейнера logging_api..."
docker run -d \
  --name logging_api \
  --restart unless-stopped \
  --network host \
  --env-file .env \
  logging_api:latest

if [ $? -eq 0 ]; then
    echo "✓ Контейнер запущен!"
    echo ""
    echo "Полезные команды:"
    echo "  Логи:    docker logs -f logging_api"
    echo "  API:     http://localhost:8080"
    echo "  Swagger: http://localhost:8080/swagger/index.html"
else
    echo "✗ Ошибка запуска контейнера"
    exit 1
fi

