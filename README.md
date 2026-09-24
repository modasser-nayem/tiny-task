# Tiny Task API

A Go-based REST API built with Gin, PostgreSQL (`pgx`), and `golang-migrate`.

---

## 🚀 Getting Started

### 1. Prerequisites
- **Go** (1.22+ or newer)
- **PostgreSQL** installed and running

### 2. Clone and Setup Environment
Clone the repository and copy the example environment file:

```bash
git clone <repository-url>
cd tiny-task
cp .env.example .env
```

Open `.env` and update your PostgreSQL credentials if needed:
```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/tiny_task?sslmode=disable
```

Make sure the database (e.g. `tiny_task`) exists in your PostgreSQL server:
```sql
CREATE DATABASE tiny_task;
```

### 3. Install Dependencies
```bash
go mod tidy
```

---

## 🗄️ Database Migrations

This project includes a built-in migration CLI at `cmd/migrate` powered by `golang-migrate`. **You do not need to install any external migration tools.**

### When Cloning for the First Time (Applying Migrations)
To create all tables and apply all pending migrations:

```bash
go run ./cmd/migrate up
```

### Check Current Migration Status
To see the active schema version:

```bash
go run ./cmd/migrate version
```

### Rollback / Undo a Migration
To rollback the most recently applied migration (1 step backward):

```bash
go run ./cmd/migrate down
```

---

## 🛠️ How to Update Tables (Adding New Migrations)

> **Golden Rule of Migrations:**  
> **NEVER edit already applied migration files** (e.g. `000001_*` or `000002_*`). Once a migration is applied or committed to Git, any schema change must be done by creating a **new migration file pair**.

### Step 1: Create New Migration Files
Inside the `migrations/` directory, create two files using the next sequential number (or timestamp):

- `000003_<description>.up.sql`
- `000003_<description>.down.sql`

#### Example: Adding a `priority` column to the `todos` table

1. Create `migrations/000003_add_priority_to_todos.up.sql`:
   ```sql
   ALTER TABLE todos ADD COLUMN priority VARCHAR(20) NOT NULL DEFAULT 'medium';
   ```

2. Create `migrations/000003_add_priority_to_todos.down.sql`:
   ```sql
   ALTER TABLE todos DROP COLUMN IF EXISTS priority;
   ```

### Step 2: Apply the New Migration
Run:
```bash
go run ./cmd/migrate up
```

Verify the version:
```bash
go run ./cmd/migrate version
```

---

## ⚠️ Troubleshooting: Fixing "Dirty" Database State

If a migration fails mid-way (for example, due to a syntax error in your SQL), `golang-migrate` locks the database in a **dirty** state to prevent data corruption.

To fix it:
1. Fix the error in your SQL file.
2. Check your current version:
   ```bash
   go run ./cmd/migrate version
   ```
3. Force the version to the last known good version (e.g. version `2`):
   ```bash
   go run ./cmd/migrate force 2
   ```
4. Re-run the migration:
   ```bash
   go run ./cmd/migrate up
   ```

---

## 🏃 Running the Server

Start the API server:

```bash
go run ./cmd/server
```

Health check:
```bash
curl http://localhost:8080/health
```

Register user endpoint:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "password123"}'
```
