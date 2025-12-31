# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a modular Golang REST API boilerplate using Fiber v2 framework with a feature-based architecture. The codebase follows a clean layered architecture pattern with clear separation between shared components and feature modules.

## Build and Run Commands

```bash
# Run the application (development)
go run cmd/api/main.go

# Build the binary
go build -o bin/api cmd/api/main.go

# Run tests
go test ./... -v

# Run tests for specific package
go test ./internal/modules/user -v

# Run single test
go test ./internal/modules/user -run TestGetProfile -v

# Run tests with coverage
go test ./... -cover

# Install dependencies
go mod download
go mod tidy
```

## Architecture

### Directory Structure

```
cmd/api/main.go          # Application entry point
internal/
  shared/                # Shared components used across modules
    config/              # Configuration loading (Viper + .env)
    database/            # Database connection (GORM + PostgreSQL)
    middleware/          # Global middleware (auth, logger, CORS, validator)
    utils/               # Utility functions (JWT, hash, response, logger, validator)
  modules/               # Feature modules
    auth/                # Authentication (login, register, refresh tokens)
    user/                # User management (CRUD)
    email/               # Email service (gomail)
    oauth/               # OAuth2 integration (Google, GitHub)
```

### Module Pattern

Each feature module follows this consistent structure:

```
module-name/
  model.go         # GORM database entity
  repository.go    # Data access layer (interface + implementation)
  service.go       # Business logic layer (interface + implementation)
  handler.go       # HTTP request handlers
  routes.go        # Route registration with middleware
  dto/
    request.go     # Input validation DTOs
    response.go    # Output DTOs
```

### Layer Responsibilities

- **model.go**: GORM entity with struct tags, relationships, and hooks
- **repository.go**: CRUD operations, database queries only (no business logic)
- **service.go**: Business logic, orchestrates repositories, transforms data
- **handler.go**: HTTP parsing, calls service, formats responses
- **routes.go**: Registers routes, applies middleware, dependency injection
- **dto/request.go**: Input structs with validation tags
- **dto/response.go**: Output structs, hides sensitive fields

### Dependency Injection Flow

The application bootstraps in `cmd/api/main.go`:

1. Load config (`config.LoadConfig()`)
2. Initialize logger
3. Initialize database connection
4. Run auto-migrations
5. Create Fiber app
6. Register global middleware (logger, CORS, recover)
7. Register module routes (each module receives `db`, `cfg`, `logger`)
8. Start server with graceful shutdown

Each module's `RegisterRoutes()` function creates its own dependency chain:
- Repository → Service → Handler → Routes

### Request Flow

```
HTTP Request → Global Middleware → Route Middleware → Handler → Service → Repository → Database
  ↓
Logger → CORS → JWT Auth → Body Validator → Parse/Validate → Business Logic → Query → Response
```

### Middleware Usage

- **BodyValidator**: Validates request against DTO struct (stores validated body in `c.Locals("validatedBody")`)
- **JWTAuth**: Protects routes by validating JWT tokens from `Authorization` header
- **HTTPLogger**: Logs all HTTP requests/responses
- **CORS**: Handles cross-origin requests

### Shared Components

**Config** (`internal/shared/config/config.go`)
- Loads from `.env` file using godotenv
- Struct with nested configs: Server, Database, JWT, OAuth, Email, Logger
- Provides `GetDSN()` for PostgreSQL connection string
- Validates required fields based on environment mode

**Database** (`internal/shared/database/connection.go`)
- GORM with PostgreSQL driver
- Connection pooling: MaxIdleConns=10, MaxOpenConns=100
- Auto-migration support via `AutoMigrate()`
- Graceful connection closing

**Utils**:
- `jwt.go`: Generate and validate JWT tokens
- `hash.go`: Password hashing with bcrypt
- `response.go`: Standardized JSON response format
- `validator.go`: Struct validation wrapper around go-playground/validator
- `logger.go`: Logrus initialization with config-based level/format

## Configuration

Copy `.env.example` to `.env` and configure:

- **SERVER_PORT**: HTTP port (default: 3000)
- **SERVER_MODE**: development/production/test
- **DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME**: PostgreSQL connection
- **JWT_SECRET**: Secret for token signing (required in production)
- **JWT_ACCESS_EXPIRY**: Access token duration (default: 1h)
- **JWT_REFRESH_EXPIRY**: Refresh token duration (default: 24h)
- **OAUTH_GOOGLE_CLIENT_ID/SECRET**: Google OAuth credentials
- **SMTP_HOST/PORT/USER/PASSWORD**: Email configuration

## Adding a New Module

1. Create module directory: `internal/modules/newmodule/dto`
2. Create files following the module pattern
3. Implement interfaces with constructors (`NewRepository`, `NewService`, `NewHandler`)
4. Create `RegisterRoutes()` function
5. In `cmd/api/main.go`: import and call `newModule.RegisterRoutes(app, db, cfg, logger)`
6. Add migrations if needed: include model in `migrationModels` slice

## Key Conventions

- **Interfaces**: Named with `I` suffix (e.g., `UserService`, `UserRepository`)
- **Implementations**: Private structs (e.g., `userService`) with `New*()` constructors
- **Repository methods**: `FindByID`, `FindAll`, `Create`, `Update`, `Delete`
- **Service methods**: Business-specific names (`GetProfile`, `CreateUser`)
- **Handler methods**: HTTP verb-based (`GetUser`, `CreateUser`)
- **Response format**: Always use `{"success": bool, "data": ..., "error": ...}` via `utils.SendResponse()`
- **Validation**: Use struct tags (`validate:"required,email,min=6"`)
- **UUID**: All entities use UUID primary keys

## Technology Stack

- **Framework**: Fiber v2
- **ORM**: GORM + PostgreSQL
- **Validation**: go-playground/validator/v10
- **JWT**: golang-jwt/jwt/v5
- **Config**: spf13/viper + joho/godotenv
- **Logger**: sirupsen/logrus
- **Email**: gopkg.in/gomail.v2
- **OAuth**: golang.org/x/oauth2
- **Testing**: stretchr/testify

## Current Modules

- **auth**: `/api/v1/auth/*` (register, login, refresh, logout)
- **user**: `/api/v1/users/*` (CRUD, requires JWT)
- **oauth**: OAuth callbacks for Google/GitHub login
- **email**: Email sending service (used by auth module)

## Notes

- All user routes except `/api/v1/auth/*` require JWT authentication
- Email module has no repository (calls external SMTP service)
- Config automatically uses default JWT secret in development mode
- Migrations run automatically on startup via `database.AutoMigrate()`
- Static files can be served from `public/` directory
