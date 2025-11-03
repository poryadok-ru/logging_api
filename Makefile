.PHONY: help build run test docker-build docker-run docker-stop docker-logs clean swagger env-example env-check start stop deploy

help: 
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

env-example:
	@echo "Создание .env.example..."
	@echo "DB_PASSWORD=your_secure_password_here" > .env.example
	@echo "Файл .env.example создан!"

env-check: 
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

docker-run: env-check ## Запустить контейнер напрямую (без compose)
	@echo "Запуск контейнера..."
	@./start.sh

docker-stop: ## Остановить контейнер
	@echo "Остановка контейнера..."
	@./stop.sh

docker-logs: ## Показать логи контейнера
	@docker logs -f logging_api

start: docker-run ## Алиас для docker-run

stop: docker-stop ## Алиас для docker-stop

clean: ## Очистить build артефакты
	@echo "Очистка..."
	rm -f logging_api
	docker stop logging_api 2>/dev/null || true
	docker rm logging_api 2>/dev/null || true
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
	@./start.sh

