# Spark Park Cricket Deployment Makefile
# Provides convenient commands for deployment and management

.PHONY: help setup-ssh deploy-dev deploy-prod deploy-all-dev deploy-all-prod build-dev build-prod test-dev test-prod health-check-dev health-check-prod clean-dev clean-prod logs-dev logs-prod rollback-dev rollback-prod

# Default target
help:
	@echo "Spark Park Cricket Deployment Commands"
	@echo "====================================="
	@echo ""
	@echo "Setup Commands:"
	@echo "  setup-ssh          Setup SSH keys for production deployment"
	@echo ""
	@echo "Deployment Commands:"
	@echo "  deploy-dev         Deploy to development environment"
	@echo "  deploy-prod        Deploy to production environment"
	@echo "  deploy-all-dev     Deploy complete stack to development"
	@echo "  deploy-all-prod    Deploy complete stack to production"
	@echo ""
	@echo "Build Commands:"
	@echo "  build-dev          Build images for development"
	@echo "  build-prod         Build images for production"
	@echo ""
	@echo "Test Commands:"
	@echo "  test-dev           Run tests in development environment"
	@echo "  test-prod          Run tests in production environment"
	@echo ""
	@echo "Health Check Commands:"
	@echo "  health-check-dev   Check health of development environment"
	@echo "  health-check-prod  Check health of production environment"
	@echo ""
	@echo "Management Commands:"
	@echo "  logs-dev           Show logs from development environment"
	@echo "  logs-prod          Show logs from production environment"
	@echo "  clean-dev          Clean up development environment"
	@echo "  clean-prod         Clean up production environment"
	@echo ""
	@echo "Rollback Commands:"
	@echo "  rollback-dev       Rollback development environment"
	@echo "  rollback-prod      Rollback production environment"
	@echo ""
	@echo "Usage Examples:"
	@echo "  make setup-ssh"
	@echo "  make deploy-prod"
	@echo "  make health-check-prod"

# Setup Commands
setup-ssh:
	@echo "Setting up SSH keys for production deployment..."
	./scripts/setup-ssh.sh --help

# Development Commands
deploy-dev:
	@echo "Deploying to development environment..."
	./scripts/deploy.sh dev --build --restart

deploy-prod:
	@echo "Deploying to production environment..."
	./scripts/deploy.sh prod --build --restart

build-dev:
	@echo "Building images for development..."
	./scripts/build.sh dev --clean

build-prod:
	@echo "Building images for production..."
	./scripts/build.sh prod --clean

deploy-all-dev:
	@echo "Deploying complete stack to development..."
	./deploy-all.sh deploy-all dev

deploy-all-prod:
	@echo "Deploying complete stack to production..."
	./deploy-all.sh deploy-all prod --with-logs

test-dev:
	@echo "Running tests in development environment..."
	cd ../backend && go test ./... -v

test-prod:
	@echo "Running tests in production environment..."
	@echo "Note: Tests are run during deployment process"

health-check-dev:
	@echo "Checking health of development environment..."
	./scripts/health-check.sh dev

health-check-prod:
	@echo "Checking health of production environment..."
	./scripts/health-check.sh prod

logs-dev:
	@echo "Showing logs from development environment..."
	docker-compose -f environments/dev/docker-compose.yml logs -f --tail=100

logs-prod:
	@echo "Showing logs from production environment..."
	ssh spark-cricket-prod "cd /opt/spark-cricket && docker-compose -f deploy/environments/prod/docker-compose.yml logs -f --tail=100"

clean-dev:
	@echo "Cleaning up development environment..."
	docker-compose -f environments/dev/docker-compose.yml down --volumes --remove-orphans
	docker system prune -f

clean-prod:
	@echo "Cleaning up production environment..."
	ssh spark-cricket-prod "cd /opt/spark-cricket && docker-compose -f deploy/environments/prod/docker-compose.yml down --volumes --remove-orphans && docker system prune -f"

rollback-dev:
	@echo "Rolling back development environment..."
	./scripts/rollback.sh dev --help

rollback-prod:
	@echo "Rolling back production environment..."
	./scripts/rollback.sh prod --help

# Utility Commands
status:
	@echo "Checking deployment status..."
	@echo ""
	@echo "Development Environment:"
	@docker-compose -f environments/dev/docker-compose.yml ps 2>/dev/null || echo "  Not running"
	@echo ""
	@echo "Production Environment:"
	@ssh spark-cricket-prod "cd /opt/spark-cricket && docker-compose -f deploy/environments/prod/docker-compose.yml ps" 2>/dev/null || echo "  Cannot connect to production server"

quick-deploy:
	@echo "Quick deployment to production (skip tests)..."
	./scripts/deploy.sh prod --build --restart --skip-tests

quick-dev:
	@echo "Quick deployment to development..."
	./scripts/deploy.sh dev --build --restart --skip-tests

# Backup Commands
backup-prod:
	@echo "Creating backup of production environment..."
	@timestamp=$$(date +%Y%m%d-%H%M%S) && \
	ssh spark-cricket-prod "cd /opt/spark-cricket && tar -czf /tmp/spark-cricket-backup-$$timestamp.tar.gz ." && \
	echo "Backup created: /tmp/spark-cricket-backup-$$timestamp.tar.gz"

restore-prod:
	@echo "Restore production environment from backup..."
	@echo "Usage: make restore-prod BACKUP_FILE=/path/to/backup.tar.gz"
	@if [ -z "$(BACKUP_FILE)" ]; then echo "Error: BACKUP_FILE not specified"; exit 1; fi
	ssh spark-cricket-prod "cd /opt/spark-cricket && tar -xzf $(BACKUP_FILE)"

# Application Commands
app-prod:
	@echo "Opening production application..."
	@echo "Application: http://15.235.202.148:3000"
	@echo "API: http://15.235.202.148:8080"

# Security Commands
update-prod:
	@echo "Updating production server packages..."
	ssh spark-cricket-prod "sudo apt update && sudo apt upgrade -y"

secure-prod:
	@echo "Applying security hardening to production server..."
	ssh spark-cricket-prod "sudo ufw enable && sudo systemctl enable fail2ban"

# Development Helpers
dev-shell:
	@echo "Opening shell in development backend container..."
	docker exec -it cricket-backend-dev /bin/sh

prod-shell:
	@echo "Opening shell in production backend container..."
	ssh spark-cricket-prod "docker exec -it cricket-backend-prod /bin/sh"

# Database Commands
db-migrate-dev:
	@echo "Running database migrations in development..."
	docker exec cricket-backend-dev ./main migrate

db-migrate-prod:
	@echo "Running database migrations in production..."
	ssh spark-cricket-prod "docker exec cricket-backend-prod ./main migrate"

# Performance Commands
benchmark-dev:
	@echo "Running performance benchmark on development..."
	cd ../backend && go test -bench=. -benchmem ./tests/performance/

benchmark-prod:
	@echo "Running performance benchmark on production..."
	@echo "Note: This will run load tests against the production server"
	@echo "Use with caution!"
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ] && \
	cd ../backend && go test -bench=. -benchmem ./tests/performance/ || echo "Cancelled"
