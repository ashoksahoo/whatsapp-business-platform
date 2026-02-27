# Setup Guide

Complete guide to setting up and running Vibecoded WA Client.

## Prerequisites

### System Requirements

- **Go**: 1.21 or higher
- **Node.js**: 18 or higher (for frontend)
- **Git**: For cloning the repository
- **Docker** (optional): For containerized deployment

### WhatsApp Business API Requirements

You need a WhatsApp Business API account with:
- WhatsApp Business Account ID
- Phone Number ID
- Access Token (API Key)
- Webhook Verify Token (you create this)

**Don't have these yet?**
1. Go to [Meta for Developers](https://developers.facebook.com/)
2. Create an app with WhatsApp product
3. Follow the [WhatsApp Business API Getting Started Guide](https://developers.facebook.com/docs/whatsapp/cloud-api/get-started)

---

## Quick Start (Development)

### 1. Clone the Repository

```bash
git clone https://github.com/ashoksahoo/whatsapp-business-platform.git
cd whatsapp-business-platform
```

### 2. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost
ENV=development

# Database (SQLite for development - no setup needed!)
DB_DRIVER=sqlite
DB_SQLITE_PATH=./whatsapp_platform.db

# WhatsApp Business Cloud API (REQUIRED)
WHATSAPP_PHONE_NUMBER_ID=your_phone_number_id_here
WHATSAPP_BUSINESS_ACCOUNT_ID=your_business_account_id_here
WHATSAPP_ACCESS_TOKEN=your_access_token_here
WHATSAPP_WEBHOOK_VERIFY_TOKEN=your_custom_verify_token_here
WHATSAPP_API_VERSION=v18.0

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run the Server

```bash
# Option 1: Direct run
go run cmd/server/main.go

# Option 2: Build and run
go build -o server cmd/server/main.go
./server

# Option 3: Using make
make run
```

**Server will start on:** `http://localhost:8080`

### 5. Verify It's Running

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2025-11-21T12:00:00Z"
}
```

---

## WhatsApp Business API Setup

### Step 1: Get Your Credentials

1. **Go to Meta for Developers**
   - Visit https://developers.facebook.com/
   - Create an app or use existing one
   - Add WhatsApp product

2. **Get Phone Number ID**
   - Go to WhatsApp > API Setup
   - Copy your Phone Number ID
   - Paste into `.env` as `WHATSAPP_PHONE_NUMBER_ID`

3. **Get Access Token**
   - In API Setup, find "Temporary Access Token"
   - Copy the token
   - Paste into `.env` as `WHATSAPP_ACCESS_TOKEN`

   **Note:** Temporary tokens expire in 24 hours. For production, generate a permanent token.

4. **Get Business Account ID**
   - Go to WhatsApp > API Setup
   - Copy your WhatsApp Business Account ID
   - Paste into `.env` as `WHATSAPP_BUSINESS_ACCOUNT_ID`

### Step 2: Configure Webhook

The webhook allows you to receive incoming messages and status updates.

1. **Make your server publicly accessible**

   For local development, use ngrok:
   ```bash
   ngrok http 8080
   ```

   Copy the HTTPS URL (e.g., `https://abc123.ngrok.io`)

2. **Set Webhook in Meta Dashboard**
   - Go to WhatsApp > Configuration
   - Click "Edit" under Webhook
   - Callback URL: `https://your-domain.com/webhooks/whatsapp`
   - Verify Token: Use the same token from `.env` (`WHATSAPP_WEBHOOK_VERIFY_TOKEN`)
   - Click "Verify and Save"

3. **Subscribe to Events**
   - Subscribe to: `messages`
   - This allows you to receive incoming messages

### Step 3: Test Sending a Message

Generate an API key first:
```sql
-- Connect to your database and run:
INSERT INTO api_keys (id, name, key_hash, permissions, created_at, updated_at)
VALUES (
  'key_test123',
  'Test API Key',
  '$2a$10$...',  -- Use bcrypt to hash your desired key
  '["messages.send", "messages.read"]',
  NOW(),
  NOW()
);
```

Then test:
```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "to": "+1234567890",
    "message_type": "text",
    "content": "Hello from Vibecoded!"
  }'
```

---

## Database Setup

### SQLite (Development - Default)

**No setup needed!** The database file is created automatically when you start the server.

- Database file: `./whatsapp_platform.db`
- Migrations run automatically on startup
- Perfect for local development and testing

**To reset the database:**
```bash
rm whatsapp_platform.db
# Restart the server - fresh database will be created
```

### PostgreSQL (Production)

For production deployments with higher load:

#### Option 1: Local PostgreSQL

1. **Install PostgreSQL**
   ```bash
   # macOS
   brew install postgresql
   brew services start postgresql

   # Ubuntu/Debian
   sudo apt install postgresql
   sudo systemctl start postgresql
   ```

2. **Create Database**
   ```bash
   createdb whatsapp_platform
   ```

3. **Update .env**
   ```env
   DB_DRIVER=postgres
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=whatsapp_platform
   DB_SSL_MODE=disable
   ```

#### Option 2: Docker PostgreSQL

```bash
docker-compose --profile postgres up -d
```

This starts PostgreSQL on port 5432 with the credentials from `.env`

---

## Docker Deployment

### Using Docker Compose (Recommended)

1. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your WhatsApp credentials
   ```

2. **Start services**
   ```bash
   # With SQLite (development)
   docker-compose up -d

   # With PostgreSQL (production)
   docker-compose --profile postgres up -d
   ```

3. **View logs**
   ```bash
   docker-compose logs -f app
   ```

4. **Stop services**
   ```bash
   docker-compose down
   ```

### Building Docker Image

```bash
# Build
docker build -t whatsapp-business-platform .

# Run
docker run -p 8080:8080 \
  --env-file .env \
  whatsapp-business-platform
```

---

## Frontend Setup

The frontend is a separate React application.

### 1. Navigate to Frontend

```bash
cd frontend
```

### 2. Install Dependencies

```bash
npm install
```

### 3. Configure Environment

```bash
cp .env.example .env
```

Edit `frontend/.env`:
```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_API_KEY=your_api_key_here
```

### 4. Run Development Server

```bash
npm run dev
```

Frontend will be available at: `http://localhost:5173`

### 5. Build for Production

```bash
npm run build
```

Built files will be in `frontend/dist/`

---

## Creating API Keys

API keys are required to use the REST API.

### Method 1: Direct Database Insert

```sql
-- Generate a key hash using bcrypt (cost 10)
-- For example, if your desired key is "my-secret-key-123"
-- Hash it using bcrypt: $2a$10$N9qo8uLOickgx2ZMRZoMye

INSERT INTO api_keys (
  id,
  name,
  key_hash,
  permissions,
  created_at,
  updated_at
) VALUES (
  'key_' || substr(md5(random()::text), 0, 16),
  'Production API Key',
  '$2a$10$yourBcryptHashHere',
  '["messages.send", "messages.read", "contacts.read", "contacts.write", "templates.read", "templates.write"]',
  NOW(),
  NOW()
);
```

### Method 2: Using Go (Recommended)

Create a helper script:

```go
// cmd/create-key/main.go
package main

import (
    "fmt"
    "github.com/ashok/whatsapp-business-platform/pkg/utils"
)

func main() {
    rawKey := "my-secret-api-key-123"
    hash, _ := utils.HashAPIKey(rawKey)
    fmt.Printf("Raw Key: %s\n", rawKey)
    fmt.Printf("Hash: %s\n", hash)
    fmt.Println("\nInsert this hash into the database")
}
```

Run:
```bash
go run cmd/create-key/main.go
```

---

## Troubleshooting

### Server Won't Start

**Error: Database connection failed**
- Check database credentials in `.env`
- Ensure PostgreSQL is running (if using postgres)
- For SQLite, check file permissions for `./whatsapp_platform.db`

**Error: Port already in use**
```bash
# Check what's using port 8080
lsof -i :8080

# Change port in .env
SERVER_PORT=8081
```

### WhatsApp Messages Not Sending

**Error: Invalid access token**
- Verify your access token hasn't expired
- Generate a new token from Meta dashboard
- Ensure token is correctly set in `.env`

**Error: Invalid phone number**
- Phone numbers must be in E.164 format: `+12345678900`
- Include country code
- No spaces or special characters

### Webhooks Not Working

**Messages not being received**
1. Verify webhook URL is publicly accessible
2. Check webhook is verified in Meta dashboard
3. Ensure you've subscribed to `messages` events
4. Check server logs: `docker-compose logs -f app`

**Webhook verification failing**
- Verify token in `.env` matches Meta dashboard
- Check webhook URL format: `https://domain.com/webhooks/whatsapp`
- Ensure server is running before verifying

### Database Issues

**SQLite: database locked**
- Only one connection can write at a time
- For high concurrency, switch to PostgreSQL

**PostgreSQL: too many connections**
- Adjust connection pool settings in code
- Check other apps aren't using all connections

---

## Environment Variables Reference

### Server
- `SERVER_PORT` - HTTP server port (default: 8080)
- `SERVER_HOST` - Host to bind to (default: localhost)
- `ENV` - Environment: development, staging, production

### Database
- `DB_DRIVER` - Database driver: sqlite or postgres
- `DB_SQLITE_PATH` - SQLite database file path
- `DB_HOST` - PostgreSQL host
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_USER` - PostgreSQL username
- `DB_PASSWORD` - PostgreSQL password
- `DB_NAME` - PostgreSQL database name
- `DB_SSL_MODE` - PostgreSQL SSL mode (disable, require)

### WhatsApp
- `WHATSAPP_PHONE_NUMBER_ID` - Your WhatsApp phone number ID (required)
- `WHATSAPP_BUSINESS_ACCOUNT_ID` - Your business account ID (required)
- `WHATSAPP_ACCESS_TOKEN` - Your API access token (required)
- `WHATSAPP_WEBHOOK_VERIFY_TOKEN` - Custom token for webhook verification (required)
- `WHATSAPP_WEBHOOK_SECRET` - Secret for webhook signature verification
- `WHATSAPP_API_VERSION` - API version (default: v18.0)

### Logging
- `LOG_LEVEL` - Log level: debug, info, warn, error
- `LOG_FORMAT` - Log format: json or console

### Storage
- `STORAGE_TYPE` - Storage type: local, s3, minio
- `MEDIA_STORAGE_PATH` - Local media storage path
- `RECORDINGS_STORAGE_PATH` - Local recordings path

---

## Production Deployment Checklist

- [ ] Use PostgreSQL instead of SQLite
- [ ] Set `ENV=production` in `.env`
- [ ] Generate permanent WhatsApp access token
- [ ] Use secure API keys (long, random strings)
- [ ] Enable SSL/TLS with reverse proxy (nginx/Traefik)
- [ ] Configure proper webhook URL (HTTPS required)
- [ ] Set up log aggregation
- [ ] Configure automated backups
- [ ] Set up monitoring and alerts
- [ ] Use environment variables (not .env file) in production
- [ ] Enable rate limiting
- [ ] Configure firewall rules

---

## Next Steps

1. **Test the API** - Use the [API Reference](./API.md) to test endpoints
2. **Run the Frontend** - Start the React frontend to manage messages
3. **Set up Webhooks** - Configure to receive incoming messages
4. **Create Templates** - Set up message templates in Meta dashboard
5. **Monitor Logs** - Watch server logs for any issues

---

## Getting Help

- **Documentation**: Check [docs/](../) folder for detailed guides
- **API Reference**: See [API.md](./API.md) for endpoint documentation
- **Issues**: This is a solo-maintained project without support channels

---

**Ready to start?** Follow the Quick Start guide above and you'll be up and running in minutes! 🚀
