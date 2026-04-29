# Paskihub Backend (paskihub-be)

Paskihub is a backend application built with **Go** (Golang) designed to manage Paskibra-related events, participants, assessments, and wallets.

## Tech Stack

- **Framework:** [Fiber v2](https://gofiber.io/)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** PostgreSQL
- **Caching:** Redis
- **Documentation:** Swagger (via `swaggo/swag`)
- **Config Management:** Viper
- **Development Tool:** Air (for live reloading)

## Project Architecture

The project follows a modular, feature-based architecture combined with some Clean Architecture principles:

- `cmd/app/main.go`: Application entry point. Initializes the server, database connection, and mounts routes.
- `domain/`: Contains core business logic definitions.
    - `contracts/`: Interfaces for repositories and services.
    - `entity/`: GORM database models.
    - `dto/`: Data Transfer Objects for requests and responses.
    - `enums/`: Shared constants and enums.
- `internal/`: Implementation details.
    - `app/`: Feature-specific modules (e.g., `user`, `event`, `assessment`). Each module typically contains:
        - `controller/`: Request handling and validation.
        - `service/`: Business logic implementation.
        - `repository/`: Data access logic (GORM).
    - `infra/`: Infrastructure components like database connection (`database/`), environment configuration (`env/`), and HTTP server setup (`server/`).
    - `middlewares/`: Custom Fiber middlewares (Auth, CORS, Rate Limiter, etc.).
- `pkg/`: Reusable utility packages (JWT, Bcrypt, UUID, Logging, etc.).
- `docs/`: Auto-generated Swagger documentation.

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL
- Redis
- [Air](https://github.com/cosmtrek/air) (optional, for development)

### Development

To run the application with live reloading:

```bash
air
```

### Building and Running

To build the application:

```bash
go build -o ./tmp/main ./cmd/app/main.go
```

To run the binary:

```bash
./tmp/main
```

### Database Migrations & Seeding

The application handles migrations automatically on startup. You can use flags to control the behavior:

- **Fresh Migration (Drop and Recreate):**
  ```bash
  go run cmd/app/main.go -fresh
  ```
- **Run Seeders:**
  ```bash
  go run cmd/app/main.go -seeder
  ```
- **Run Specific Seeder:**
  ```bash
  go run cmd/app/main.go -seeder -model <ModelName>
  ```

## Development Conventions

### 1. Feature Organization
When adding a new feature, create a new directory in `internal/app/` with `controller`, `service`, and `repository` subdirectories. Define the corresponding contracts in `domain/contracts/` and entities in `domain/entity/`.

### 2. Error Handling
Use the custom error handling defined in `domain/errors.go` or standard Go error patterns. Ensure errors are logged appropriately using `pkg/log`.

### 3. API Documentation
Update Swagger comments in controllers and run `swag init` from the root directory to refresh documentation. Swagger UI is available at `/swagger/` in development mode.

### 4. Validation
Use `github.com/go-playground/validator/v10` for request validation. Custom validators can be registered in `internal/infra/server/http_server.go`.

### 5. Authentication
Authentication is handled via JWT. Use the `middlewares.Auth()` middleware to protect routes. Role-based access control should be checked within services or via specialized middlewares.

## Key Files

- `cmd/app/main.go`: Entry point and wiring.
- `.air.toml`: Air configuration for live reload.
- `internal/infra/server/http_server.go`: Route definitions and server configuration.
- `internal/infra/database/pgsql_conn.go`: Database connection and migration logic.
- `go.mod`: Dependency management.
