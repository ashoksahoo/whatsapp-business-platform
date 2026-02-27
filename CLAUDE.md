# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Vibecoded WA Client is a self-hosted WhatsApp Business API client with Go backend and React frontend. It provides REST API access to WhatsApp Business Cloud API with features for messaging, contact management, templates, and webhooks.

## Build and Development Commands

### Backend (Go)

```bash
# Run the server (development)
go run cmd/server/main.go

# Or using make
make run

# Build the binary
make build

# Run tests
go test -v ./...
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt
go fmt ./...

# Lint (requires golangci-lint)
make lint

# Download dependencies
go mod download
go mod tidy
# Or using make
make deps

# Development with hot reload (requires air)
make dev

# Install development tools
make install-tools
```

### Frontend (React + Vite)

```bash
cd frontend

# Install dependencies
npm install

# Run development server (port 5173)
npm run dev

# Build for production
npm run build

# Run linter
npm run lint

# Preview production build
npm run preview
```

### Docker

```bash
# Start with SQLite (default)
docker-compose up

# Start with PostgreSQL (production profile)
docker-compose --profile postgres up

# Stop all services
make docker-down

# View logs
make docker-logs
```

## Architecture

### Layered Architecture

The backend follows a clean layered architecture:

1. **API Layer** (`internal/api/`)
   - `handlers/` - HTTP request handlers
   - `middleware/` - Auth, CORS, logging, rate limiting, recovery, request ID
   - `routes/` - Route registration and grouping
   - `server.go` - Server initialization and lifecycle

2. **Service Layer** (`internal/services/`)
   - Business logic and orchestration
   - Services: MessageService, ContactService, TemplateService, AuthService
   - Communicates with repositories and external APIs (WhatsApp)

3. **Repository Layer** (`internal/repositories/`)
   - Data access abstraction
   - GORM-based implementations
   - Repositories: MessageRepository, ContactRepository, TemplateRepository, APIKeyRepository

4. **Model Layer** (`internal/models/`)
   - Database models with GORM tags
   - Base model with UUID, timestamps, soft delete
   - Custom types (MessageType, MessageStatus, Direction, MediaType)

5. **Integration Layer** (`internal/whatsapp/`)
   - WhatsApp Business Cloud API client
   - Retry logic, error handling, structured logging
   - Methods: SendTextMessage, SendMediaMessage, SendTemplateMessage, GetMessageStatus

### Dependency Flow

```
HTTP Request → Handler → Service → Repository → Database
                    ↓
                WhatsApp Client → WhatsApp API
```

Initialization (in `cmd/server/main.go`):
1. Load config from environment/`.env`
2. Initialize logger (zap)
3. Connect to database (SQLite or PostgreSQL)
4. Run auto-migrations
5. Initialize: Repositories → Services → Handlers → Routes
6. Start HTTP server with graceful shutdown

### Database

- **Development**: SQLite (default, `./vibecoded.db`)
- **Production**: PostgreSQL (configurable via `DB_DRIVER`)
- **ORM**: GORM with auto-migrations
- **Features**: UUID primary keys, soft deletes, timestamps, indexes, triggers

Switch database by setting `DB_DRIVER=postgres` or `DB_DRIVER=sqlite` in `.env`

### Configuration

Configuration via Viper from `.env` file and environment variables:
- Server (port, host, environment, shutdown timeout)
- Database (driver, connection params, pool settings)
- WhatsApp (API token, phone number ID, webhook tokens)
- Security (API key salt, session secret)
- Logging (level, format, output)
- Storage (local/S3/Minio for media and recordings)
- Metrics (Prometheus port)

Validation enforces required fields in production mode.

## API Structure

All authenticated endpoints require `X-API-Key` header and use `/api/v1` prefix:

- **Messages**: `POST /messages`, `GET /messages`, `GET /messages/search`, `GET /messages/:id`
- **Contacts**: `GET /contacts`, `GET /contacts/search`, `GET /contacts/:id`, `PATCH /contacts/:id`
- **Templates**: `GET /templates`, `POST /templates`, `GET /templates/:id`, `PATCH /templates/:id`, `DELETE /templates/:id`
- **Webhooks**: `GET /webhooks/whatsapp` (verification), `POST /webhooks/whatsapp` (receive)
- **Health**: `GET /health` (no auth)

Rate limiting: 1000 requests/minute per API key

## Frontend Architecture

React 18 + TypeScript + Vite + Tailwind CSS

**Design System**: "Strawberry Theme"
- Primary: `#f43f5e` (rose-500)
- Success: `#22c55e` (green-500)
- Font: Inter
- Components use Tailwind classes with strawberry color palette

**Pages**: Messages, Contacts, Templates

**State Management**: React hooks (no Redux)

**HTTP Client**: Axios configured with base URL and API key from `.env`

## Testing

### Backend
```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./internal/services

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestFunctionName ./path/to/package
```

### Frontend
```bash
cd frontend
npm run lint
```

## Environment Setup

Copy `.env.example` to `.env` and configure:

**Required for development**:
- `WHATSAPP_ACCESS_TOKEN` - From Meta Business
- `WHATSAPP_PHONE_NUMBER_ID` - From Meta Business
- `WHATSAPP_WEBHOOK_VERIFY_TOKEN` - Your custom token

**Database** (defaults to SQLite):
```env
DB_DRIVER=sqlite
DB_SQLITE_PATH=./vibecoded.db
```

For PostgreSQL in production:
```env
DB_DRIVER=postgres
DB_HOST=your-host
DB_PORT=5432
DB_USER=your-user
DB_PASSWORD=your-password
DB_NAME=vibecoded_wa
```

## Key Implementation Details

### WhatsApp Message Flow

**Outbound**: Handler → MessageService → WhatsApp Client → Meta API → Database (saved)

**Inbound**: Meta Webhook → WebhookHandler (signature verification) → MessageService → Database

### Authentication

API key based authentication via `X-API-Key` header. Keys are hashed with salt before storage. AuthMiddleware validates all `/api/v1/*` routes.

### Error Handling

Custom error types in `pkg/errors/`:
- ValidationError
- NotFoundError
- WhatsAppError
- InternalError

Middleware captures panics and returns 500 with request ID for debugging.

### Logging

Structured logging with zap:
- JSON format in production
- Console format in development
- Request ID attached to all logs
- Different log levels per environment

### Database Models

All models extend BaseModel (UUID, CreatedAt, UpdatedAt, DeletedAt). Key models:
- **Message**: WhatsApp messages (text, media, template) with status tracking
- **Contact**: WhatsApp contacts (auto-created from messages)
- **Template**: Message templates with parameters
- **APIKey**: Hashed API keys with metadata

## Module Path Note

The go.mod uses `github.com/ashok/vibecoded-wa-client` as the module path. All internal imports use this path.

There is a local replacement for `github.com/square-key-labs/strawgo-ai` pointing to `../strawgo-ai` (for local development).

## Coming Features

Voice calling implementation is planned with SIP integration, WebRTC, call recording, and AI transcription (Whisper/Deepgram). Related files are stubbed in `internal/voice/` and `internal/models/call.go`.
