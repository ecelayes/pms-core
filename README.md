# 🏨 Hotel PMS Core (Backend)

Scalable backend for a Property Management System (PMS) and Booking Engine.
Built with **Go (Golang)** following **Clean Architecture** and **Zero Repetition** principles.

## 🚀 Tech Stack

* **Language:** Go 1.21+
* **Web Framework:** Echo v4
* **Database:** PostgreSQL 14+
* **SQL Driver:** pgx/v5 (High performance)
* **Migrations:** golang-migrate
* **Architecture:** Hexagonal / Clean Architecture

---

## 🛠️ Prerequisites

Before starting, ensure you have the following installed:

1. **Go**: [Download](https://go.dev/dl/)
2. **PostgreSQL**: [Download](https://www.postgresql.org/download/)
3. **Migrate CLI**: Tool to manage database migrations.

```bash
# Install Migrate CLI
go install -tags 'postgres' [github.com/golang-migrate/migrate/v4/cmd/migrate@latest](https://github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
```

---

## ⚡ Quick Start (How to run the project)

Follow these steps to set up the development environment:

### 1. Database Configuration

Create an empty database in your local Postgres engine:

```sql
CREATE DATABASE hotel_pms_db;
```

### 2. Environment Variables

Create a `.env` file in the project root (you can use the example below):

```ini
# .env
PORT=8081
DATABASE_URL=postgres://postgres:postgres@localhost:5432/hotel_pms_db?sslmode=disable
```

### 3. Run Migrations

Use the `Makefile` to create the tables automatically:

```bash
make db-up
```

### 4. Seed Data (Optional)

Populate the database with dummy data (Hotels, Rooms, Inventory) for testing:

```bash
make db-seed
```

### 5. Start the Server

Start the application:

```bash
make run
```

*You should see:* `Starting server on port 8081...`

---

## 📦 Makefile Commands

The project includes a `Makefile` to automate repetitive tasks. No need to memorize long commands.

| Command | Description | Usage Example |
| :--- | :--- | :--- |
| `make run` | Starts the server reading variables from `.env`. | `make run` or `make run PORT=3000` |
| `make db-up` | Applies all pending migrations (creates tables). | `make db-up` |
| `make db-down` | Reverts the last applied migration (drops tables). | `make db-down` |
| `make db-create` | Creates a new empty migration `.sql` file. | `make db-create name=add_users_table` |
| `make db-force` | Fixes the DB if a migration failed and state is "dirty". | `make db-force version=1` |
| `make test` | Runs all unit tests. | `make test` |

---

## 📂 Project Structure

We follow the **Standard Go Project Layout** adapted to Clean Architecture:

```text
pms-core/
├── cmd/api/                # Entry Point (Main)
├── internal/               # Private application code
│   ├── core/               # Pure Business Logic
│   │   ├── domain/         # Entities (Structs without DB logic)
│   │   └── ports/          # Interfaces (Contracts for Repos/Services)
│   ├── adapter/            # Concrete Implementations
│   │   ├── handler/        # HTTP Controllers (Echo)
│   │   └── repository/     # Data Access (Postgres/pgx)
│   └── infra/              # Tooling Config (Server, DB connect)
├── migrations/             # SQL Files (Up/Down)
├── pkg/                    # Shared Code/Utils (Responses, Loggers)
└── Makefile                # Automation Scripts
```

---

## 🔌 Available Endpoints

### Public

* `GET /api/v1/health`: Service health check.
* `GET /api/v1/availability`: Check room availability.
  * **Params:** `start` (YYYY-MM-DD), `end` (YYYY-MM-DD)
  * **Example:** `http://localhost:8081/api/v1/availability?start=2023-12-01&end=2023-12-05`

---

## 🗄️ Data Model

The system features the following main tables (see `migrations/` for details):

1. **Hotels & Users:** Multi-tenant configuration.
2. **Inventory:** Central table optimized for fast availability reads.
3. **Bookings:** Booking transactions.
4. **Extras:** Additional services (Breakfast, Parking, etc.).
