#!/bin/bash

echo "Остановка контейнера logging_api..."
docker stop logging_api

echo "Удаление контейнера..."
docker rm logging_api

echo "Готово!"

