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
    database/            # Database connection (GORM + PostgreSQL) + migrations + redis
    middleware/          # Global middleware (auth, logger, CORS, validator, RBAC)
    utils/               # Utility functions (JWT, hash, response, logger, validator, random)
  modules/               # Feature modules
    auth/                # Authentication (login, register, refresh tokens, verification)
    user/                # User management (CRUD)
    role/                # Role and permission management (RBAC)
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
- **service.go**: Business logic, orchestrates repositories, transforms data, integrates third party services (Email, Redis)
- **handler.go**: HTTP parsing, calls service, formats responses
- **routes.go**: Registers routes, applies middleware, dependency injection
- **dto/request.go**: Input structs with validation tags
- **dto/response.go**: Output structs, hides sensitive fields

### Dependency Injection Flow

The application bootstraps in `cmd/api/main.go`:

1. Load config (`config.LoadConfig()`)
2. Initialize logger
3. Initialize database connection (PostgreSQL)
4. Initialize Redis connection
5. Run migrations (manual via `cmd/migrate` or auto in dev)
6. Seed initial roles (SuperAdmin, Admin, User)
7. Create Fiber app
8. Register global middleware (logger, CORS, recover)
9. Register module routes (each module receives `db`, `cfg`, `logger`, `redisClient`)
10. Start server with graceful shutdown

Each module's `RegisterRoutes()` function creates its own dependency chain:
- Repository → Service → Handler → Routes

### Request Flow

```
HTTP Request → Global Middleware → Route Middleware → Handler → Service → Repository → Database
  ↓                                                                     ↓
Logger → CORS → JWT Auth → RBAC Check → Body Validator → Parse/Validate → Business Logic → Query → Response
                                                                        ↓
                                                                      Redis (OTP/Cache)
                                                                        ↓
                                                                      Email (SMTP)
```

### Middleware Usage

- **BodyValidator**: Validates request against DTO struct (stores validated body in `c.Locals("validatedBody")`)
- **JWTAuth**: Protects routes by validating JWT tokens from `Authorization` header
- **RequireRole**: Checks if authenticated user has any of the specified roles (admin, super_admin)
- **RequirePermission**: Checks if authenticated user has a specific permission (users.create, roles.update)
- **HTTPLogger**: Logs all HTTP requests/responses
- **CORS**: Handles cross-origin requests

## Security Features

The API supports several security features that can be enabled/disabled via environment variables:

### Account Activation (Email Verification)
- **Flag**: `EMAIL_VERIFICATION_ENABLED`
- **Flow**: New users receive a 6-digit OTP via email and must verify it before they can log in.
- **Exceptions**: SuperAdmin is automatically verified.

### Two-Factor Authentication (2FA)
- **Flag**: `TWO_FACTOR_ENABLED`
- **Flow**: After entering password, users receive a 6-digit OTP via email and must provide it to receive tokens.
- **Exceptions**: SuperAdmin is exempt from 2FA flow.

### Session Management
- **Flow**: Refresh tokens are stored in the database as **Sessions** with device metadata.
- **Metadata Recorded**: IP Address, User Agent, Device ID (from `X-Device-ID` header).
- **Features**: List active sessions, logout from specific devices, block specific sessions.

## Database & Migrations

The project uses `golang-migrate` for versioned migrations.

- **Migrations Path**: `db/migrations/`
- **CLI Tool**: `go run cmd/migrate/main.go`
- **Commands**: `-up`, `-down`, `-create NAME`

### Table Naming Convention

Tables use prefixes to indicate their type:

**Master Tables** (prefix `m_`):
- `m_users` - User accounts
- `m_roles` - Role definitions

**Transaction Tables** (prefix `t_`):
- `t_sessions` - User sessions and refresh tokens (contains device metadata)
- `t_oauth_accounts` - OAuth provider links

## RBAC System (Role-Based Access Control)

### Overview

The API implements a comprehensive RBAC system with:
- **3 Default Roles**: SuperAdmin, Admin, User
- **Granular Permissions**: Format `resource.action` (e.g., `users.create`, `roles.delete`)
- **Wildcard Permission**: `*` grants full access (SuperAdmin only)
- **Role Storage**: Separate `m_roles` table with foreign key to `m_users`
- **Stateless Auth**: Role and permission data embedded in JWT tokens
- **JSONB Storage**: Permissions stored as JSONB type in PostgreSQL for efficient querying

### Default Roles and Permissions

**SuperAdmin** (`slug: super_admin`)
- Permissions: `["*"]` (full access to everything)
- Can: Manage all resources, assign roles, manage roles

**Admin** (`slug: admin`)
- Permissions: `["users.create", "users.read", "users.update", "users.delete", "roles.read", "roles.assign"]`
- Can: Create/read/update/delete users, read roles, assign roles to users
- Cannot: Manage roles (create/update/delete roles)

**User** (`slug: user`)
- Permissions: `["users.read", "users.update"]` (own profile only)
- Can: Read and update own profile
- Cannot: Access other users, manage roles, perform admin operations

### Using RBAC Middleware

**RequireRole - Protect routes by role:**
```go
// Only SuperAdmin can access
protected.Use(middleware.RequireRole(cfg, "super_admin"))

// Admin or SuperAdmin can access
protected.Use(middleware.RequireRole(cfg, "admin", "super_admin"))
```

**RequirePermission - Protect routes by permission:**
```go
// Only users with users.create permission can access
protected.Use(middleware.RequirePermission(cfg, "users.create"))
```

**Helper Functions:**
```go
// Get user role from context
roleSlug, ok := middleware.GetRoleSlugFromContext(c)

// Get user permissions from context
permissions, ok := middleware.GetPermissionsFromContext(c)

// Get user ID from context
userID, ok := middleware.GetUserIDFromContext(c)
```

### JWT Claims Structure

JWT tokens include role information:
```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "role_slug": "admin",
  "permissions": ["users.create", "users.read", "users.update"],
  "exp": 1234567890
}
```

### Role Assignment Rules

The API enforces strict role assignment rules to maintain security:

**Registration (POST /api/v1/auth/register):**
- Automatically assigns "user" role
- Cannot specify role during registration
- All new users start with basic "user" permissions

**Create User (POST /api/v1/users) - Admin/SuperAdmin only:**
- Can optionally specify `role_id` in request body
- Only allows creating users with "user" or "admin" roles
- Cannot create users with "super_admin" role via this endpoint
- If `role_id` is not provided, defaults to "user" role
- Example: `{"name": "John", "email": "john@example.com", "password": "pass123", "role_id": "uuid-here"}`

**Update User (PUT /api/v1/users/:id):**
- Admin/SuperAdmin can update `role_id` field
- Only allows updating role to "user" or "admin"
- Regular users cannot update their own role (blocked at handler level)
- Non-admin users can only update their name and email

**Assign Role (PUT /api/v1/users/:id/role) - SuperAdmin only:**
- Can assign any role including "super_admin"
- This is the ONLY way to grant super_admin role to a user
- Requires role UUID in request body

**Summary Table:**

| Endpoint | Access Level | Can Assign "user"? | Can Assign "admin"? | Can Assign "super_admin"? |
|----------|--------------|-------------------|---------------------|---------------------------|
| **POST /api/v1/auth/register** | Public | ✅ (auto) | ❌ | ❌ |
| **POST /api/v1/users** | Admin/SuperAdmin | ✅ (default) | ✅ (optional) | ❌ (blocked) |
| **PUT /api/v1/users/:id** | All users* | ✅ (admin only) | ✅ (admin only) | ❌ (blocked) |
| **PUT /api/v1/users/:id/role** | SuperAdmin only | ✅ | ✅ | ✅ |

*Regular users can update their own profile but NOT their role. Only Admin/SuperAdmin can update roles.

### Protected Routes Summary

**Public Routes:**
- `/api/v1/auth/register` - User registration
- `/api/v1/auth/login` - User login
- `/api/v1/auth/refresh` - Token refresh
- `/api/v1/oauth/*` - OAuth redirects and callbacks

**Authenticated Routes (Any User):**
- `/api/v1/users/me` - Get/update own profile
- `/api/v1/users/:id` (PUT) - Update user (self or admin)
- `/api/v1/auth/sessions` (GET) - List all active sessions
- `/api/v1/auth/sessions/:id` (DELETE) - Logout from a specific device
- `/api/v1/auth/sessions/:id/block` (PATCH) - Block a specific session

**Admin/SuperAdmin Routes:**
- `/api/v1/users` (GET) - List all users
- `/api/v1/users` (POST) - Create user
- `/api/v1/users/:id` (DELETE) - Delete user
- `/api/v1/roles` (GET) - List all roles

**SuperAdmin Only Routes:**
- `/api/v1/users/:id/role` (PUT) - Assign role to user
- `/api/v1/roles` (POST) - Create role
- `/api/v1/roles/:id` (PUT/DELETE) - Update/delete role

## Database Table Naming Convention

Tables use prefixes to indicate their type:

**Master Tables** (prefix `m_`):
- `m_users` - User accounts
- `m_roles` - Role definitions

**Transaction Tables** (prefix `t_`):
- `t_sessions` - User sessions and refresh tokens
- `t_oauth_accounts` - OAuth provider links

**Migration Strategy:**
- In development mode, old tables (`users`, `oauth_accounts`, `refresh_tokens`) are dropped on startup
- New tables with prefixes are created automatically via GORM AutoMigrate
- This is controlled by the `RenameTables()` function in `internal/shared/database/migration.go`
- Only runs when `SERVER_MODE=development`

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
- **SUPERADMIN_NAME**: Default SuperAdmin account name (default: "Super Admin")
- **SUPERADMIN_EMAIL**: Default SuperAdmin email (default: "superadmin@boilerplate.com")
- **SUPERADMIN_PASSWORD**: Default SuperAdmin password (default: "SuperAdmin123!")

### SuperAdmin Account

The application automatically creates/updates a default SuperAdmin account on startup using credentials from `.env`:
- If the account doesn't exist, it will be created
- If the account exists, password and details will be updated from `.env` config
- Always assigned the "super_admin" role with full `["*"]` permissions
- **Important**: Change the default password after first login in production!

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
- **user**: `/api/v1/users/*` (CRUD with role-based access control)
- **role**: `/api/v1/roles/*` (role management, SuperAdmin only)
- **oauth**: `/api/v1/oauth/*` (Google/GitHub OAuth)
- **email**: Email sending service (used by auth and oauth modules)

## Notes

- All user routes except `/api/v1/auth/*` require JWT authentication
- RBAC middleware enforces role and permission-based access control
- Email module has no repository (calls external SMTP service)
- Config automatically uses default JWT secret in development mode
- Migrations run automatically on startup via `database.AutoMigrate()`
- Initial roles are seeded automatically on first startup
- Table rename migration runs in development mode to drop old tables
- Static files can be served from `public/` directory

## Feature Flags

Optional features can be enabled/disabled via environment variables:

- **OAUTH_GOOGLE_ENABLED**: Enable/disable Google OAuth (default: false)
- **OAUTH_GOOGLE_SEND_WELCOME_EMAIL**: Send welcome email after Google OAuth (default: false)
- **OAUTH_GITHUB_ENABLED**: Enable/disable GitHub OAuth (default: false)
- **OAUTH_GITHUB_SEND_WELCOME_EMAIL**: Send welcome email after GitHub OAuth (default: false)
- **EMAIL_ENABLED**: Master switch for email functionality (default: false)
