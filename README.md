# WhatsApp Business Platform

> Modern, open-source WhatsApp Business Platform

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)](https://react.dev/)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg)](CONTRIBUTING.md)

**[Website](https://ashok.io/whatsapp-business-platform/)** · **[API Docs](./docs/API.md)** · **[Setup Guide](./docs/SETUP.md)**

## What is This?

A production-ready, self-hosted WhatsApp Business API platform. Built with a Go backend and React frontend, it gives you full control over your WhatsApp Business integration — no third-party SaaS required.

## Features

### Backend (REST API)
- Send text messages via WhatsApp
- Send media messages (images, documents, audio, video)
- Send template messages
- Receive incoming messages via webhooks
- Contact management (auto-creation, search, updates)
- Message history and search
- Template management
- API key authentication
- Rate limiting (1000 req/min per key)
- Structured logging with zap

### Frontend (React + Tailwind)
- Messages page with real-time messaging interface
- Contacts page with search and management
- Templates page for message templates
- Responsive layout with toast notifications
- TypeScript for type safety

### Infrastructure
- SQLite (development) / PostgreSQL (production)
- Docker & Docker Compose setup
- Environment-based configuration
- Health check endpoints
- Auto-migrations
- Graceful shutdown

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- WhatsApp Business API credentials
- Docker (optional)

### Backend Setup

1. **Clone the repository**
```bash
git clone https://github.com/ashoksahoo/whatsapp-business-platform.git
cd whatsapp-business-platform
```

2. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your WhatsApp API credentials
```

3. **Run the backend**
```bash
# Development (uses SQLite by default)
go run cmd/server/main.go

# Or with make
make run
```

The API will be available at `http://localhost:8080`

### Frontend Setup

1. **Install dependencies**
```bash
cd frontend
npm install
```

2. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your API configuration
```

3. **Run the frontend**
```bash
npm run dev
```

The frontend will be available at `http://localhost:5173`

### Using Docker

```bash
# Start everything (uses SQLite)
docker-compose up

# With PostgreSQL (for production)
docker-compose --profile postgres up
```

## Database Options

### Development (SQLite - Default)
No setup needed. The SQLite database is created automatically.

```env
DB_DRIVER=sqlite
DB_SQLITE_PATH=./data.db
```

### Production (PostgreSQL)

```env
DB_DRIVER=postgres
DB_HOST=your-postgres-host
DB_PORT=5432
DB_USER=your-user
DB_PASSWORD=your-password
DB_NAME=whatsapp_platform
```

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin
- **ORM**: GORM
- **Database**: SQLite / PostgreSQL
- **API**: WhatsApp Business Cloud API

### Frontend
- **Framework**: React 18
- **Build Tool**: Vite
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Router**: React Router
- **HTTP Client**: Axios
- **Icons**: Lucide React

## Roadmap

### Voice Calling
- Make & receive voice calls on WhatsApp
- Automatic call recording
- AI transcription (Whisper/Deepgram)
- Voice AI agents
- Call analytics

### Enhanced Messaging
- Interactive messages (buttons, lists)
- Message scheduling
- Bulk messaging
- Analytics dashboard
- Multi-account support
- Webhook forwarding

### Integrations
- MCP server for Claude AI
- Prometheus metrics
- S3/Minio for media storage

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Documentation

- **[Setup Guide](./docs/SETUP.md)** - Complete setup and deployment guide
- **[API Reference](./docs/API.md)** - Complete REST API documentation
- [Architecture](./docs/ARCHITECTURE.md) - System architecture
- [Data Models](./docs/DATA_MODELS.md) - Database schema

## Use Cases

- **Customer Support**: Handle WhatsApp inquiries at scale
- **Notifications**: Send order updates and transactional alerts
- **Marketing**: Send promotional messages via templates
- **Integration**: Connect WhatsApp to your existing systems
- **AI Assistants**: Build chatbots and automated responses

## Environment Variables

See [.env.example](./.env.example) for all available configuration options.

Key variables:
- `DB_DRIVER` - Database driver (`sqlite` / `postgres`)
- `WHATSAPP_ACCESS_TOKEN` - Your WhatsApp API token
- `WHATSAPP_PHONE_NUMBER_ID` - Your WhatsApp phone number ID
- `SERVER_PORT` - API server port (default: `8080`)

## Troubleshooting

**Database Issues**
- SQLite: Check file permissions for the DB path
- PostgreSQL: Verify connection credentials in `.env`

**WhatsApp API Issues**
- Verify your access token is valid and not expired
- Check webhook URL is publicly accessible
- Ensure phone number ID is correct

**Frontend Connection**
- Verify backend is running on the correct port
- Check CORS settings if accessing from a different origin
- Ensure API key is set in frontend `.env`

## License

MIT License — See [LICENSE](LICENSE) for details.

---

**If you find this useful, please star the repo!**
