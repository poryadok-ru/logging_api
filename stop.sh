#!/bin/bash

echo "Остановка контейнера logging_api..."
docker stop logging_api 2>/dev/null

if [ $? -eq 0 ]; then
    echo "Удаление контейнера..."
    docker rm logging_api 2>/dev/null
    echo "✓ Контейнер остановлен и удалён"
else
    echo "⚠️  Контейнер не найден или уже остановлен"
fi

