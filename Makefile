.PHONY: help build run test docker-build docker-run docker-stop docker-logs clean swagger deploy lint

help: ## Показать список доступных команд
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

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

docker-run: ## Запустить контейнер
	@echo "Запуск контейнера..."
	@./start.sh

docker-stop: ## Остановить контейнер
	@echo "Остановка контейнера..."
	@./stop.sh

docker-logs: ## Показать логи контейнера
	@docker logs -f logging_api

clean: ## Очистить build артефакты и контейнеры
	@echo "Очистка..."
	@rm -f logging_api
	@docker stop logging_api 2>/dev/null || true
	@docker rm logging_api 2>/dev/null || true
	@docker rmi logging_api:latest 2>/dev/null || true
	@echo "✓ Очистка завершена"

deploy: docker-build docker-run ## Полный деплой (сборка + запуск)

lint: ## Запустить линтер
	@echo "Запуск линтера..."
	golangci-lint run

