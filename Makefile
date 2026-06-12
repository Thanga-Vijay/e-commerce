# E-Commerce Platform Makefile
# Simplify Docker operations for development and production

.PHONY: help build up down logs clean restart rebuild ps test lint

# Default target
.DEFAULT_GOAL := help

# Colors for terminal output
BLUE := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

help: ## Show this help message
	@echo "$(BLUE)E-Commerce Platform - Docker Management$(RESET)"
	@echo ""
	@echo "$(GREEN)Available commands:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(RESET) %s\n", $$1, $$2}'

# Development Commands
dev-build: ## Build all services for development
	@echo "$(BLUE)Building development environment...$(RESET)"
	docker-compose build

dev-up: ## Start development environment
	@echo "$(GREEN)Starting development environment...$(RESET)"
	docker-compose up -d
	@echo "$(GREEN)Services are running!$(RESET)"
	@echo "Frontend: http://localhost:3000"
	@echo "Auth Service: http://localhost:8081"

dev-up-build: ## Build and start development environment
	@echo "$(GREEN)Building and starting development environment...$(RESET)"
	docker-compose up -d --build

dev-down: ## Stop development environment
	@echo "$(RED)Stopping development environment...$(RESET)"
	docker-compose down

dev-logs: ## View logs from all development services
	docker-compose logs -f

dev-logs-service: ## View logs from specific service (make dev-logs-service SERVICE=auth-service)
	docker-compose logs -f $(SERVICE)

dev-restart: ## Restart development environment
	@echo "$(YELLOW)Restarting development environment...$(RESET)"
	docker-compose restart

dev-ps: ## List running development containers
	docker-compose ps

# Production Commands
prod-build: ## Build all services for production
	@echo "$(BLUE)Building production environment...$(RESET)"
	docker-compose -f docker-compose.prod.yml build

prod-up: ## Start production environment
	@echo "$(GREEN)Starting production environment...$(RESET)"
	docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d
	@echo "$(GREEN)Production services are running!$(RESET)"

prod-up-build: ## Build and start production environment
	@echo "$(GREEN)Building and starting production environment...$(RESET)"
	docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d --build

prod-down: ## Stop production environment
	@echo "$(RED)Stopping production environment...$(RESET)"
	docker-compose -f docker-compose.prod.yml down

prod-logs: ## View logs from all production services
	docker-compose -f docker-compose.prod.yml logs -f

prod-ps: ## List running production containers
	docker-compose -f docker-compose.prod.yml ps

# Staging Commands
staging-build: ## Build all services for staging
	@echo "$(BLUE)Building staging environment...$(RESET)"
	docker-compose -f docker-compose.prod.yml build

staging-up: ## Start staging environment
	@echo "$(GREEN)Starting staging environment...$(RESET)"
	docker-compose -f docker-compose.prod.yml --env-file .env.staging up -d
	@echo "$(GREEN)Staging services are running!$(RESET)"

staging-down: ## Stop staging environment
	@echo "$(RED)Stopping staging environment...$(RESET)"
	docker-compose -f docker-compose.prod.yml down

# Database Commands
db-migrate: ## Run database migrations for all services
	@echo "$(BLUE)Running database migrations...$(RESET)"
	@echo "Auth Service..."
	docker-compose exec auth-service migrate -path migrations -database "postgresql://postgres:postgres@auth-db:5432/auth_db?sslmode=disable" up || true
	@echo "Product Service..."
	docker-compose exec product-service migrate -path migrations -database "postgresql://postgres:postgres@product-db:5432/product_db?sslmode=disable" up || true
	@echo "$(GREEN)Migrations complete!$(RESET)"

db-seed: ## Seed databases with sample data
	@echo "$(BLUE)Seeding databases...$(RESET)"
	# Add your seed commands here
	@echo "$(GREEN)Seeding complete!$(RESET)"

db-backup: ## Backup all databases
	@echo "$(BLUE)Backing up databases...$(RESET)"
	mkdir -p backups
	docker-compose exec -T auth-db pg_dump -U postgres auth_db > backups/auth_db_$$(date +%Y%m%d_%H%M%S).sql
	docker-compose exec -T product-db pg_dump -U postgres product_db > backups/product_db_$$(date +%Y%m%d_%H%M%S).sql
	docker-compose exec -T cart-db pg_dump -U postgres cart_db > backups/cart_db_$$(date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)Backups complete in ./backups/$(RESET)"

db-restore: ## Restore databases from backup (make db-restore FILE=backup.sql DB=auth_db)
	@echo "$(BLUE)Restoring database $(DB)...$(RESET)"
	docker-compose exec -T $(DB)-db psql -U postgres $(DB) < $(FILE)
	@echo "$(GREEN)Restore complete!$(RESET)"

# Maintenance Commands
clean: ## Remove all stopped containers, volumes, and images
	@echo "$(RED)Cleaning up Docker resources...$(RESET)"
	docker-compose down -v
	docker system prune -af --volumes
	@echo "$(GREEN)Cleanup complete!$(RESET)"

clean-volumes: ## Remove all volumes (WARNING: deletes all data)
	@echo "$(RED)Removing all volumes...$(RESET)"
	docker-compose down -v
	@echo "$(GREEN)Volumes removed!$(RESET)"

rebuild: ## Rebuild and restart all services
	@echo "$(YELLOW)Rebuilding all services...$(RESET)"
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d
	@echo "$(GREEN)Rebuild complete!$(RESET)"

restart-service: ## Restart specific service (make restart-service SERVICE=auth-service)
	@echo "$(YELLOW)Restarting $(SERVICE)...$(RESET)"
	docker-compose restart $(SERVICE)

# Monitoring Commands
stats: ## Show resource usage statistics
	docker stats

health: ## Check health of all services
	@echo "$(BLUE)Checking service health...$(RESET)"
	@docker-compose ps | grep -E "(healthy|Up)"

inspect-service: ## Inspect specific service (make inspect-service SERVICE=auth-service)
	docker-compose exec $(SERVICE) sh

# Image Commands
pull: ## Pull latest images
	@echo "$(BLUE)Pulling latest images...$(RESET)"
	docker-compose pull

push: ## Push images to registry (requires Docker Hub login)
	@echo "$(BLUE)Pushing images to registry...$(RESET)"
	docker-compose -f docker-compose.prod.yml build
	docker-compose -f docker-compose.prod.yml push

tag: ## Tag images for registry (make tag VERSION=1.0.0)
	@echo "$(BLUE)Tagging images with version $(VERSION)...$(RESET)"
	docker tag ecommerce/auth-service:latest ecommerce/auth-service:$(VERSION)
	docker tag ecommerce/product-service:latest ecommerce/product-service:$(VERSION)
	docker tag ecommerce/cart-service:latest ecommerce/cart-service:$(VERSION)
	docker tag ecommerce/wishlist-service:latest ecommerce/wishlist-service:$(VERSION)
	docker tag ecommerce/order-service:latest ecommerce/order-service:$(VERSION)
	docker tag ecommerce/payment-service:latest ecommerce/payment-service:$(VERSION)
	docker tag ecommerce/inventory-service:latest ecommerce/inventory-service:$(VERSION)
	docker tag ecommerce/notification-service:latest ecommerce/notification-service:$(VERSION)
	docker tag ecommerce/reporting-service:latest ecommerce/reporting-service:$(VERSION)
	docker tag ecommerce/frontend:latest ecommerce/frontend:$(VERSION)
	@echo "$(GREEN)Tagging complete!$(RESET)"

# Testing Commands
test: ## Run tests for all services
	@echo "$(BLUE)Running tests...$(RESET)"
	# Add your test commands here
	@echo "$(GREEN)Tests complete!$(RESET)"

test-service: ## Run tests for specific service (make test-service SERVICE=auth-service)
	@echo "$(BLUE)Running tests for $(SERVICE)...$(RESET)"
	docker-compose exec $(SERVICE) go test ./...

# Security Commands
security-scan: ## Run security scan on all images
	@echo "$(BLUE)Running security scan...$(RESET)"
	docker scan ecommerce/auth-service:latest || true
	docker scan ecommerce/product-service:latest || true
	@echo "$(GREEN)Security scan complete!$(RESET)"

# Network Commands
network-inspect: ## Inspect Docker networks
	docker network ls
	docker network inspect ecommerce-network

# Quick Commands
quick-start: dev-up ## Quick start development environment
	@echo "$(GREEN)Development environment started!$(RESET)"
	@echo "Visit http://localhost:3000"

quick-stop: dev-down ## Quick stop development environment
	@echo "$(RED)Development environment stopped!$(RESET)"

quick-restart: dev-down dev-up ## Quick restart development environment
	@echo "$(GREEN)Development environment restarted!$(RESET)"

# Information Commands
version: ## Show Docker and Docker Compose versions
	@echo "$(BLUE)Docker version:$(RESET)"
	@docker --version
	@echo "$(BLUE)Docker Compose version:$(RESET)"
	@docker-compose --version

info: ## Show system information
	@echo "$(BLUE)Docker System Info:$(RESET)"
	@docker info | grep -E "Server Version|Operating System|Total Memory|CPUs"
	@echo ""
	@echo "$(BLUE)Running Containers:$(RESET)"
	@docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

ports: ## Show exposed ports
	@echo "$(BLUE)Service Ports:$(RESET)"
	@echo "Frontend:      http://localhost:3000"
	@echo "Auth Service:  http://localhost:8081"
	@echo "Product:       http://localhost:8082"
	@echo "Cart:          http://localhost:8083"
	@echo "Wishlist:      http://localhost:8084"
	@echo "Order:         http://localhost:8085"
	@echo "Payment:       http://localhost:8086"
	@echo "Inventory:     http://localhost:8087"
	@echo "Notification:  http://localhost:8088"
	@echo "Reporting:     http://localhost:8089"

# Kafka Commands (Phase 9)
kafka-up: ## Start Kafka infrastructure
	@echo "$(GREEN)Starting Kafka infrastructure...$(RESET)"
	docker-compose -f docker-compose.kafka.yml up -d
	@echo "$(GREEN)Kafka is running!$(RESET)"
	@echo "Kafka UI: http://localhost:8090"

kafka-down: ## Stop Kafka infrastructure
	@echo "$(RED)Stopping Kafka infrastructure...$(RESET)"
	docker-compose -f docker-compose.kafka.yml down

kafka-logs: ## View Kafka logs
	docker-compose -f docker-compose.kafka.yml logs -f kafka

kafka-topics: ## List all Kafka topics
	@echo "$(BLUE)Kafka Topics:$(RESET)"
	docker-compose -f docker-compose.kafka.yml exec kafka kafka-topics --list --bootstrap-server localhost:9092

kafka-topic-describe: ## Describe specific topic (make kafka-topic-describe TOPIC=order.created)
	docker-compose -f docker-compose.kafka.yml exec kafka kafka-topics --describe --topic $(TOPIC) --bootstrap-server localhost:9092

kafka-consume: ## Consume messages from topic (make kafka-consume TOPIC=order.created)
	docker-compose -f docker-compose.kafka.yml exec kafka kafka-console-consumer --topic $(TOPIC) --from-beginning --bootstrap-server localhost:9092

kafka-produce: ## Produce message to topic (make kafka-produce TOPIC=order.created)
	docker-compose -f docker-compose.kafka.yml exec kafka kafka-console-producer --topic $(TOPIC) --bootstrap-server localhost:9092

kafka-consumer-groups: ## List all consumer groups
	docker-compose -f docker-compose.kafka.yml exec kafka kafka-consumer-groups --list --bootstrap-server localhost:9092

kafka-consumer-lag: ## Check consumer group lag (make kafka-consumer-lag GROUP=your-group)
	docker-compose -f docker-compose.kafka.yml exec kafka kafka-consumer-groups --describe --group $(GROUP) --bootstrap-server localhost:9092

kafka-reset-offset: ## Reset consumer offset (make kafka-reset-offset GROUP=your-group TOPIC=your-topic)
	docker-compose -f docker-compose.kafka.yml exec kafka kafka-consumer-groups --reset-offsets --to-earliest --execute --group $(GROUP) --topic $(TOPIC) --bootstrap-server localhost:9092

kafka-ui: ## Open Kafka UI in browser
	@echo "$(BLUE)Opening Kafka UI...$(RESET)"
	@echo "Kafka UI: http://localhost:8090"

# Full Stack Commands
stack-up: ## Start full stack (services + Kafka + monitoring)
	@echo "$(GREEN)Starting full stack...$(RESET)"
	docker-compose up -d
	docker-compose -f docker-compose.kafka.yml up -d
	docker-compose -f docker-compose.monitoring.yml up -d
	@echo "$(GREEN)Full stack is running!$(RESET)"
	@echo "Frontend: http://localhost:3000"
	@echo "Kafka UI: http://localhost:8090"
	@echo "Grafana: http://localhost:3001"
	@echo "Prometheus: http://localhost:9090"

stack-down: ## Stop full stack
	@echo "$(RED)Stopping full stack...$(RESET)"
	docker-compose down
	docker-compose -f docker-compose.kafka.yml down
	docker-compose -f docker-compose.monitoring.yml down

stack-logs: ## View logs from full stack
	docker-compose logs -f & docker-compose -f docker-compose.kafka.yml logs -f

monitoring-up: ## Start monitoring stack
	@echo "$(GREEN)Starting monitoring stack...$(RESET)"
	docker-compose -f docker-compose.monitoring.yml up -d
	@echo "$(GREEN)Monitoring is running!$(RESET)"
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://localhost:3001 (admin/admin)"
	@echo "Alertmanager: http://localhost:9093"

monitoring-down: ## Stop monitoring stack
	@echo "$(RED)Stopping monitoring stack...$(RESET)"
	docker-compose -f docker-compose.monitoring.yml down
