# Vibecoded WA Client

> 🍓 A self-hosted WhatsApp Business API client with a beautiful React frontend

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)](https://react.dev/)

## 🎯 What is This?

A production-ready WhatsApp Business API client that you can self-host. Built with Go backend and React frontend, featuring a beautiful strawberry-themed UI.

## ✨ Available Features

### Backend (REST API)
- ✅ Send text messages via WhatsApp
- ✅ Send media messages (images, documents, audio, video)
- ✅ Send template messages
- ✅ Receive incoming messages via webhooks
- ✅ Contact management (auto-creation, search, updates)
- ✅ Message history and search
- ✅ Template management
- ✅ API key authentication
- ✅ Rate limiting
- ✅ Comprehensive logging

### Frontend (React + Tailwind)
- ✅ Messages page with real-time messaging interface
- ✅ Contacts page with search and management
- ✅ Templates page for message templates
- ✅ Strawberry theme design system
- ✅ Responsive layout
- ✅ Toast notifications
- ✅ TypeScript for type safety

### Infrastructure
- ✅ SQLite database (development) / PostgreSQL (production)
- ✅ Docker & Docker Compose setup
- ✅ Environment-based configuration
- ✅ Health check endpoints
- ✅ Auto-migrations
- ✅ Graceful shutdown

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- WhatsApp Business API credentials
- Docker (optional)

### Backend Setup

1. **Clone the repository**
```bash
git clone https://github.com/ashoksahoo/vibecode-wa-business.git
cd vibecode-wa-business
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

## 📊 Database Options

### Development (SQLite - Default)
No setup needed! The SQLite database is created automatically.

```env
DB_DRIVER=sqlite
DB_SQLITE_PATH=./vibecoded.db
```

### Production (PostgreSQL)
For production deployments with higher load:

```env
DB_DRIVER=postgres
DB_HOST=your-postgres-host
DB_PORT=5432
DB_USER=your-user
DB_PASSWORD=your-password
DB_NAME=vibecoded_wa
```

Switch between databases by changing the `DB_DRIVER` environment variable.

## 🎨 Frontend Features

Built with the **Strawberry Theme** design system:

- **Primary Color**: Strawberry Red (#f43f5e)
- **Success Color**: Leaf Green (#22c55e)
- **Font**: Inter
- **Components**: Fully responsive, accessible, with smooth animations

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin (web framework)
- **ORM**: GORM
- **Database**: SQLite (dev) / PostgreSQL (prod)
- **API**: WhatsApp Business Cloud API

### Frontend
- **Framework**: React 18
- **Build Tool**: Vite
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Router**: React Router
- **HTTP Client**: Axios
- **Icons**: Lucide React

## 📁 Project Structure

```
vibecoded-wa-client/
├── backend/
│   ├── cmd/server/          # Application entry point
│   ├── internal/
│   │   ├── api/             # HTTP handlers & routes
│   │   ├── config/          # Configuration
│   │   ├── database/        # Database connection
│   │   ├── models/          # Data models
│   │   ├── repositories/    # Data access layer
│   │   ├── services/        # Business logic
│   │   └── whatsapp/        # WhatsApp client
│   └── pkg/                 # Utilities
│
└── frontend/
    ├── src/
    │   ├── components/      # UI components
    │   ├── pages/           # Page components
    │   ├── layouts/         # Layout components
    │   ├── services/        # API client
    │   └── types/           # TypeScript types
    └── dist/                # Production build
```

## 🔜 Coming Soon

### Voice Calling 🎙️
- Make & receive voice calls on WhatsApp
- Automatic call recording
- AI transcription (Whisper/Deepgram)
- Voice AI agents
- Searchable call history
- Call analytics

### Enhanced Features
- Interactive messages (buttons, lists)
- Message scheduling
- Bulk messaging
- Analytics dashboard
- Multi-account support
- Webhook forwarding

### Integrations
- MCP server for Claude AI integration
- Prometheus metrics
- S3/Minio for media storage

## 🚫 Contribution Policy

This is a **solo-maintained project**. I am not accepting pull requests or external contributions.

**However, you are free to:**
- ✅ Fork and modify
- ✅ Use commercially (MIT License)
- ✅ Create derivative works

## 📖 Documentation

- **[API Reference](./docs/API.md)** - Complete REST API documentation
- [API Design](./docs/API_DESIGN.md) - API design specifications
- [Architecture](./docs/ARCHITECTURE.md) - System architecture
- [Data Models](./docs/DATA_MODELS.md) - Database schema
- [UI/UX Guide](./docs/design/UI_UX_GUIDE.md) - Design system
- [Frontend README](./frontend/README.md) - Frontend documentation

## 🎯 Use Cases

- **Customer Support**: Handle WhatsApp inquiries
- **Notifications**: Send order updates and alerts
- **Marketing**: Send promotional messages via templates
- **Integration**: Connect WhatsApp to your existing systems
- **AI Assistants**: Build chatbots and automated responses

## 📝 Environment Variables

See [.env.example](./.env.example) for all available configuration options.

Key variables:
- `DB_DRIVER` - Database driver (sqlite/postgres)
- `WHATSAPP_ACCESS_TOKEN` - Your WhatsApp API token
- `WHATSAPP_PHONE_NUMBER_ID` - Your WhatsApp phone number ID
- `SERVER_PORT` - API server port (default: 8080)

## 🐛 Troubleshooting

### Database Issues
- SQLite: Check file permissions for `./vibecoded.db`
- PostgreSQL: Verify connection credentials in `.env`

### WhatsApp API Issues
- Verify your access token is valid
- Check webhook URL is publicly accessible
- Ensure phone number ID is correct

### Frontend Connection
- Verify backend is running on the correct port
- Check CORS settings if accessing from different origin
- Ensure API key is set in frontend `.env`

## 📝 License

MIT License - See [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Built with [Claude](https://claude.ai) - AI-assisted development at its finest!

---

**⭐ If you find this useful, please star the repo!**

Built with 🍓 by Ashok | MIT License | Solo-maintained
