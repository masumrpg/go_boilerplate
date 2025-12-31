# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **Golang modular web application boilerplate** using Fiber v2 framework with a **feature-based modular architecture**. The project emphasizes clear separation of concerns with a strict layered structure within each module.

**Current State**: This is a newly initialized project with architectural documentation but no implementation code yet. The structure below is planned and documented in README.md (in Indonesian).

## Common Development Commands

```bash
# Build (once cmd/api/main.go exists)
go build -o bin/app ./cmd/api

# Run application
go run cmd/api/main.go

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests in specific package
go test ./internal/modules/user/...

# Run specific test function
go test -v ./internal/modules/user/... -run TestLogin

# Dependency management
go mod download
go mod tidy
go mod verify

# Static analysis
go vet ./...
gofmt -l .
```

**Note**: No Makefile or build scripts exist yet. Use standard Go commands.

## Architecture Overview

### Directory Structure

```
cmd/api/main.go          # Application entry point (initialize deps, register routes, start server)
internal/
├── shared/              # Cross-cutting concerns used by all modules
│   ├── config/         # Configuration management (Viper, environment-based)
│   ├── database/       # Database connection (GORM + PostgreSQL, connection pooling)
│   ├── middleware/     # HTTP middleware (auth, logger, CORS, validator)
│   └── utils/          # Utility functions (JWT, hash, validator, response, logger)
└── modules/            # Feature modules (self-contained business features)
    ├── auth/           # Authentication & authorization
    ├── user/           # User management
    ├── email/          # Email notifications
    └── oauth/          # Social login (Google, GitHub)
docs/                   # Swagger/OpenAPI documentation
pkg/                    # Public packages (optional, for external reuse)
```

### Architectural Pattern

- **Feature-based modules**: Each business feature is a self-contained module
- **Clean layers**: Handler → Service → Repository (strict separation)
- **Dependency injection**: Constructor-based DI in main.go
- **Interface-driven**: Each layer defines interfaces for testability
- **DTO pattern**: Request/response DTOs for API contracts and data hiding

## Module Structure (CRITICAL - Must Follow)

Every module MUST follow this exact 7-file pattern:

```
module-name/
├── model.go           # Domain entity (GORM model with tags, relationships, hooks)
├── repository.go      # Data access layer (interface + CRUD implementation)
├── service.go         # Business logic layer (interface + implementation)
├── handler.go         # HTTP request handlers (parse request, call service, format response)
├── routes.go          # Route registration (register endpoints, apply middleware, DI)
└── dto/
    ├── request.go     # Input DTOs with validation tags
    └── response.go    # Output DTOs (hide sensitive fields like password)
```

**Example**: Creating a new module
```bash
mkdir -p internal/modules/payment/dto
cd internal/modules/payment
touch model.go repository.go service.go handler.go routes.go
touch dto/request.go dto/response.go
```

## Layer Responsibilities

### model.go
- Define database entity/schema with GORM struct tags
- Table relationships (belongs to, has many, etc.)
- GORM hooks (BeforeCreate, AfterUpdate, etc.)
- **NO business logic**

### repository.go
- Interface definition for data access
- CRUD operations (Create, FindByID, Update, Delete, FindAll)
- Database queries using GORM
- **NO business logic** - only data access

### service.go
- Interface definition for business operations
- Business logic implementation
- Orchestrate multiple repositories if needed
- Call external services (email, OAuth, etc.)
- Transform domain models to DTOs
- **NO HTTP or database-specific code**

### handler.go
- Parse HTTP request and extract parameters
- Validate input using DTOs
- Call service methods
- Format and return HTTP response
- Handle HTTP-specific errors
- **NO business logic** - delegate to service

### routes.go
- Register routes for the module
- Apply middleware (JWT auth, validation, etc.)
- Group related endpoints
- Dependency injection: receive dependencies from main.go and wire them together
- Example pattern:
  ```go
  func RegisterRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
      userRepo := repository.NewUserRepository(db)
      userService := service.NewUserService(userRepo)
      userHandler := handler.NewUserHandler(userService)

      api := app.Group("/api/v1")
      api.Post("/users", userHandler.CreateUser)
      api.Get("/users/:id", middleware.AuthRequired, userHandler.GetUser)
  }
  ```

### dto/request.go
- Input validation structs with validation tags
- Request body parsing
- Example:
  ```go
  type CreateUserRequest struct {
      Name  string `json:"name" validate:"required,min=3,max=100"`
      Email string `json:"email" validate:"required,email"`
  }
  ```

### dto/response.go
- Output format structs
- Hide sensitive fields (password, tokens, etc.)
- Example:
  ```go
  type UserResponse struct {
      ID    uuid.UUID `json:"id"`
      Name  string    `json:"name"`
      Email string    `json:"email"`
      // Password is intentionally omitted
  }
  ```

## Request Flow

```
HTTP Request
    ↓
Global Middleware (Logger, CORS)
    ↓
Fiber Router (Route Matching)
    ↓
Route-specific Middleware (JWT Auth, Validator)
    ↓
Handler Layer
    ├─ Parse Request Body → DTO
    ├─ Validate Input
    └─ Call Service
         ↓
Service Layer
    ├─ Execute Business Logic
    ├─ Call Repository (or multiple repositories)
    ├─ Transform Domain Models → DTOs
    └─ Call External Services (Email, OAuth)
         ↓
Repository Layer
    ├─ Execute GORM Queries
    └─ Return Domain Models
         ↓
Database (PostgreSQL via GORM)
```

## Naming Conventions

### Files
- Use singular nouns: `model.go` (User, Product, Order)

### Interfaces and Structs
- **Interface**: `UserRepository`, `UserService` (PascalCase, exported)
- **Implementation**: `userRepository`, `userService` (camelCase, private)
- **Constructor**: `NewUserRepository()`, `NewUserService()`, `NewUserHandler()`

### Methods
- **Repository layer**: `Create`, `FindByID`, `Update`, `Delete`, `FindAll`, `FindByEmail`
- **Service layer**: Business-specific names (`GetUserProfile`, `CreateOrder`, `ProcessPayment`, `AuthenticateUser`)
- **Handler layer**: HTTP verb-like (`GetUser`, `CreateUser`, `UpdateUser`, `DeleteUser`)

## Dependency Injection Pattern

Use **constructor-based dependency injection** in `main.go`:

```go
// 1. Initialize shared components
cfg := config.Load()
db := database.Connect(cfg)
logger := logger.New(cfg)

// 2. Create Fiber app
app := fiber.New()

// 3. Register global middleware
app.Use(middleware.Logger())
app.Use(middleware.CORS())

// 4. Register module routes with DI
user.RegisterRoutes(app, db, cfg, logger)
auth.RegisterRoutes(app, db, cfg, logger)
email.RegisterRoutes(app, db, cfg, logger)

// 5. Start server
app.Listen(":3000")
```

Inside each module's `routes.go`:
```go
func RegisterRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config, logger *logrus.Logger) {
    // Wire dependencies: Repository → Service → Handler
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    // Register routes
    api := app.Group("/api/v1")
    api.Post("/users", userHandler.CreateUser)
    api.Get("/users/:id", userHandler.GetUser)
}
```

## Module Development Workflow

When creating a new module, follow this 7-step process:

1. **Create Model** - Define GORM schema in `model.go`
2. **Create Repository** - Define interface and implement CRUD operations
3. **Create Service** - Define interface and implement business logic
4. **Create DTOs** - Define request/response DTOs with validation
5. **Create Handler** - Implement HTTP handlers
6. **Create Routes** - Register endpoints and wire dependencies
7. **Register in main.go** - Call the module's RegisterRoutes function

## Shared Components

### internal/shared/config/
- Uses Viper for configuration management
- Load configuration from environment variables
- Provides type-safe config struct

### internal/shared/database/
- GORM + PostgreSQL connection pooling
- Auto migration support
- Connection management

### internal/shared/middleware/
- **auth.go**: JWT validation middleware
- **logger.go**: HTTP request/response logging
- **cors.go**: CORS configuration
- **validator.go**: Request body validation

### internal/shared/utils/
- **jwt.go**: Generate & validate JWT tokens
- **hash.go**: Password hashing with bcrypt
- **validator.go**: Struct validation helper
- **response.go**: Standardized JSON response format
- **logger.go**: Logrus configuration

## Technology Stack

- **Framework**: Fiber v2 (github.com/gofiber/fiber/v2)
- **ORM**: GORM + PostgreSQL
- **Validation**: go-playground/validator/v10
- **JWT**: golang-jwt/jwt/v5
- **Config**: spf13/viper
- **Logger**: sirupsen/logrus
- **Email**: gopkg.in/gomail.v2
- **OAuth**: golang.org/x/oauth2
- **Testing**: stretchr/testify

## Key Architectural Principles

1. **Separation of Concerns**: Each layer has a single, well-defined responsibility
2. **Module Independence**: Modules are self-contained and loosely coupled
3. **Interface-Driven Design**: Use interfaces for testability and flexibility
4. **DTO Pattern**: Never expose domain models directly to API; use DTOs
5. **Dependency Injection**: Wire dependencies explicitly in main.go
6. **No Layer Violations**: Handlers don't access database; services don't know about HTTP

## Testing Strategy

- Test each layer in isolation
- Mock dependencies using interfaces
- Test business logic in service layer
- Test HTTP handlers with test requests
- Example test file: `service_test.go`, `handler_test.go`

```go
// Example: Mock repository for service testing
type mockUserRepository struct {}
func (m *mockUserRepository) FindByID(id uuid.UUID) (*User, error) {
    // Return test data
}
```

## Important Notes

- README.md contains detailed architectural documentation in Indonesian
- Original documentation includes Mermaid diagrams showing request flow, DI flow, and module interactions
- This project uses Go 1.25.5
- All JSON serialization uses bytedance/sonic (faster than standard JSON)
- API documentation is generated using swaggo/swag
