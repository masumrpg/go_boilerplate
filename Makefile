# Start local development server
run:
	go run cmd/api/main.go

# Build binary
build:
	go build -o bin/api cmd/api/main.go

# Test
test:
	go test ./... -v

# Database Migrations
MIGRATE_CMD = go run cmd/migrate/main.go

migrate-up:
	$(MIGRATE_CMD) -up

migrate-down:
	$(MIGRATE_CMD) -down

migrate-force:
	@read -p "Enter version to force: " version; \
	$(MIGRATE_CMD) -force $$version

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir db/migrations -seq $$name

# Docker
docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Swagger
swagger:
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

# Module Generator
module:
	@read -p "Enter module name (singular): " name; \
	go run cmd/gen/main.go $$name
