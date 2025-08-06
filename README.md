# Courtly API
- Go API for managing sports court reservations, Pix payments, and sending notifications. API implementation for [Courtly](https://github.com/dinizgab/courtly-web).

## Requirements.
- Go 1.24+
- PostgreSQL
- Goose for migrations
- Docker (optional for local environment)

## Environment Variables

| Name                | Description |
|---------------------|-------------|
| API_PORT            | Port where the API will be exposed |
| JWT_SECRET          | Key used to sign JWT tokens |
| DATABASE_URL        | PostgreSQL database connection URL |
| SMTP_EMAIL          | Sender used for sending emails |
| SMTP_HOST           | SMTP server host |
| SMTP_PORT           | SMTP server port |
| SMTP_USER           | SMTP server user |
| SMTP_PASS           | SMTP server password |
| OPENPIX_BASE_URL    | Base URL for the OpenPix API |
| OPENPIX_APP_ID      | OpenPix application ID |
| STORAGE_PROJECT_URL | Supabase storage project URL |
| STORAGE_API_KEY     | API key for storage |

## Running Locally

Start support services (optional):
```bash
make compose-up
```

Run database migrations:
```bash
make migration-up
```

Start the application:
```bash
make run
```

The API will be accessible at **http://localhost:$API_PORT**.

## Structure
- `cmd/main.go` – application entry point.
- `internal/` – domain modules, repositories, use cases, and handlers implementation.
- `migrations/` – SQL scripts for database creation and modification.
