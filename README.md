# RSS Aggregator

A backend REST API for aggregating RSS feeds, built with Go. This is my first Go project created for learning purposes.

## Overview

RSS Aggregator is a simple but functional backend service that allows users to create accounts, subscribe to RSS feeds, and retrieve posts from those feeds. The application includes a background worker that periodically scrapes configured RSS feeds and stores the latest posts in a PostgreSQL database.

## Features

- User management with API key authentication
- Feed subscription and management
- Automatic RSS feed scraping via background worker
- Posts aggregation from multiple RSS feeds
- RESTful API with CORS support
- Graceful shutdown handling
- PostgreSQL database with migrations

## Prerequisites

- Go 1.25.4 or higher
- Docker & Docker Compose
- Goose CLI (for database migrations)
- HTTP CLI tool (for testing endpoints)

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/viniciuspra/rssagg.git
cd rssagg
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Environment Setup

Create a `.env` file in the project root:

```env
PORT=8080
DB_USER=docker
DB_PASSWORD=your_secure_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=rssagg
```

The application automatically constructs the database connection URL from these variables with SSL mode disabled.

### 4. Start Database

```bash
make p-up
```

This starts the PostgreSQL container using Docker Compose.

### 5. Run Migrations

Install Goose if you haven't already:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Then run migrations:

```bash
make g-up
```

To check migration status:

```bash
make g-status
```

### 6. Generate Database Code

If you modify any SQL queries, regenerate the Go code:

```bash
make sqlc
```

## Running the Application

Build and run the application:

```bash
make dev
```

Or manually:

```bash
go build
./rssagg
```

The server will start on `http://localhost:8080`.

## API Endpoints

All endpoints are prefixed with `/v1`. Protected endpoints require `Authorization: ApiKey {your_api_key}` header.

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/healthz` | No | Health check |
| POST | `/users` | No | Create user (returns API key) |
| GET | `/users` | Yes | Get current user |
| POST | `/feeds` | Yes | Create feed |
| GET | `/feeds` | No | Get all feeds |
| POST | `/feedFollows` | Yes | Subscribe to feed |
| GET | `/feedFollows` | Yes | Get user subscriptions |
| DELETE | `/feedFollows/{id}` | Yes | Unsubscribe from feed |
| GET | `/posts` | Yes | Get posts from subscribed feeds (max 10) |

## Usage Example

### 1. Create a User

```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice"}' | jq .
```

Save the returned `api_key`.

### 2. Create a Feed

```bash
curl -X POST http://localhost:8080/v1/feeds \
  -H "Authorization: ApiKey YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name": "BBC", "url": "http://feeds.bbc.co.uk/news/rss.xml"}' | jq .
```

Save the returned `id`.

### 3. Subscribe to Feed

```bash
curl -X POST http://localhost:8080/v1/feedFollows \
  -H "Authorization: ApiKey YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"feed_id": "FEED_ID_FROM_STEP_2"}' | jq .
```

### 4. View Posts

```bash
curl -X GET http://localhost:8080/v1/posts \
  -H "Authorization: ApiKey YOUR_API_KEY" | jq .
```

The background scraper runs every minute and fetches new posts from all feeds in the database.

## Database Migrations

Migrations are located in `sql/migrations/` and use Goose.

To create a new migration:

```bash
goose -dir sql/migrations create migration_name sql
```

## Development Tasks

Available Makefile commands:

```bash
make dev              # Build and run the application
make p-up             # Start PostgreSQL container
make p-stop           # Stop PostgreSQL container
make g-up             # Run all pending migrations
make g-down           # Rollback last migration
make g-status         # Show migration status
make sqlc             # Generate Go code from SQL queries
make health           # Check server health
```

## Authentication

The API uses API key authentication via the `Authorization` header:

```
Authorization: ApiKey {your_api_key}
```

API keys are generated automatically when creating a user and are stored as SHA256 hashes.

## Background Worker

The application includes a background scraper that:

- Runs in a separate goroutine
- Fetches up to 10 feeds every minute
- Parses RSS feeds and extracts post data
- Stores new posts in the database
- Skips duplicate posts (checks URL uniqueness)
- Handles feed fetch failures gracefully

## Graceful Shutdown

The application handles SIGINT and SIGTERM signals, allowing:

- Active requests to complete (up to 10 seconds timeout)
- Database connections to close properly
- Background workers to stop cleanly

Press Ctrl+C to trigger shutdown.

## Learning Notes

This project serves as a learning exercise in:

- Building REST APIs with Go
- Working with databases and migrations
- Implementing authentication middleware
- Background workers and concurrency
- Error handling and logging
- Context and lifecycle management
- Docker for development environments

## Future Improvements

Potential areas for enhancement:

- Add pagination to posts endpoint
- Implement post filtering and search
- Add feed categorization
- User preferences for update frequency
- Rate limiting
- More comprehensive error handling
- Unit and integration tests
- API documentation with Swagger

## License

This project is for educational purposes.
