.PHONY: help build run test docker-build docker-run docker-stop clean swagger env-example env-check

help: ## Показать помощь
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

env-example: ## Создать .env.example файл (шаблон)
	@echo "Создание .env.example..."
	@echo "DB_PASSWORD=your_secure_password_here" > .env.example
	@echo "Файл .env.example создан!"

env-check: ## Проверить наличие .env файла
	@if [ ! -f .env ]; then \
		echo "⚠️  Файл .env не найден!"; \
		echo "Создайте файл .env с содержимым:"; \
		echo ""; \
		echo "DB_PASSWORD=your_secure_password_here"; \
		echo ""; \
		echo "Или запустите: cp .env.example .env и отредактируйте"; \
		exit 1; \
	else \
		echo "✓ Файл .env найден"; \
	fi

build: ## Скомпилировать приложение
	@echo "Компиляция..."
	go build -o logging_api .

run: ## Запустить приложение локально
	@echo "Запуск приложения..."
	go run .

test: ## Запустить тесты
	@echo "Запуск тестов..."
	go test -v ./...

swagger: ## Сгенерировать Swagger документацию
	@echo "Генерация Swagger..."
	swag init

docker-build: ## Собрать Docker образ
	@echo "Сборка Docker образа..."
	docker build -t logging_api:latest .

docker-run: env-check ## Запустить через Docker Compose
	@echo "Запуск через Docker Compose..."
	docker-compose up -d

docker-stop: ## Остановить Docker контейнеры
	@echo "Остановка контейнеров..."
	docker-compose down

docker-logs: ## Показать логи контейнера
	docker-compose logs -f logging_api

clean: ## Очистить build артефакты
	@echo "Очистка..."
	rm -f logging_api
	docker-compose down -v
	docker rmi logging_api:latest 2>/dev/null || true

deploy: env-check docker-build docker-run ## Полный деплой (проверка .env, сборка, запуск) 

lint:
	@echo "Запуск линтера..."
	golangci-lint run

prod-build: ## Собрать для продакшена
	@echo "Сборка production образа..."
	docker build -t logging_api:prod --build-arg BUILD_ENV=production .

prod-deploy: prod-build
	@echo "Деплой на продакшен..."
	docker-compose -f docker-compose.prod.yml up -d

