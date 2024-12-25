# Указываем путь к docker-compose файлу
COMPOSE_FILE := docker-compose.yml

# Переменные окружения
ENV_FILE := .env

# Сборка Docker-образов
build:
	@echo "Building Docker images..."
	docker-compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) build

# Запуск контейнеров
up:
	@echo "Starting services..."
	docker-compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) up -d

# Остановка контейнеров
down:
	@echo "Stopping services..."
	docker-compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) down

# Перезапуск контейнеров
restart: down up

# Просмотр логов
logs:
	@echo "Displaying logs..."
	docker-compose -f $(COMPOSE_FILE) --env-file $(ENV_FILE) logs -f

# Удаление остановленных контейнеров и неиспользуемых ресурсов
cleanup:
	@echo "Cleaning up unused resources..."
	docker system prune -f

# Помощь
help:
	@echo "Usage:"
	@echo "  make build     - Build Docker images"
	@echo "  make up        - Start services"
	@echo "  make down      - Stop services"
	@echo "  make restart   - Restart services"
	@echo "  make logs      - Display logs"
	@echo "  make cleanup   - Remove unused Docker resources"
