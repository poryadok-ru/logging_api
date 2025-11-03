#!/bin/bash

API_URL="https://api.alexmayka.ru"

echo "=== Нагрузочное тестирование Logging API ==="
echo ""

echo "1. Health Check (простой эндпоинт)"
wrk -t4 -c100 -d10s --latency ${API_URL}/health
echo ""

echo "2. GET /api/v1/auth/me (с авторизацией)"
wrk -t2 -c20 -d10s --latency -s auth_token.lua ${API_URL}/api/v1/auth/me
echo ""

echo "3. GET /api/v1/bots (список ботов)"
wrk -t2 -c20 -d10s --latency -s auth_token.lua ${API_URL}/api/v1/bots
echo ""

echo "4. POST /api/v1/logs (создание логов)"
wrk -t2 -c10 -d10s --latency -s post_log.lua ${API_URL}/api/v1/logs
echo ""

echo "=== Тестирование завершено ==="

