# Project Shop API

Project Shop API — a sample backend API for managing in-game items, player inventories, and purchases.

**Overview**
- **Language:** Go
- **Purpose:** A focused API for item shop operations (item listing, stock management, purchasing, and purchase history)
- **Main folders:** Code is organized into domain packages such as `pkg/itemShop`, `pkg/inventory`, `pkg/itemManaging`, `server`, and `databases`

**Key Features**
- Create, update, and manage items
- Shop system: list items, check stock, purchase items, and record purchase history
- Player inventory management: add/remove items for players
- OAuth2 (Google) authentication for selected endpoints

**Prerequisites**
- Go 1.20+ (see `go.mod`)
- PostgreSQL (configuration in `config/config.yml`)

**Setup**
1. Copy and edit `config/config.yml` to set database and OAuth credentials
2. Create the PostgreSQL database and run migrations (migration scripts are under `databases/migration`)

**Run (example)**
```bash
go build -o shop-api .
./shop-api
```

Or run directly:
```bash
go run main.go
```

**Brief Project Structure**
- `main.go` — application entry point
- `config/` — application configuration and `config.yml`
- `databases/` — database connection and migration scripts
- `pkg/` — domain packages (itemShop, inventory, itemManaging, player, admin, logger, oauth2)
- `server/` — router and middleware setup
- `tests/` — example unit/integration tests

**Testing**
Run unit tests for the repository or specific packages:
```bash
go test ./...
```

**Suggestions / Further Improvements**
- Add package-level README files (for `pkg/itemShop`, `pkg/inventory`, etc.) describing API endpoints, request/response examples
- Provide a `Makefile` or `docker-compose.yml` for easier local development (database + app)

**Contact / Next Steps**
- If you want further edits to this README, help running tests, or committing and pushing changes, let me know.

