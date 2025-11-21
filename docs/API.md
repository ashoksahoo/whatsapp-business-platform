# API Reference

Complete REST API documentation for Vibecoded WA Client.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All API endpoints (except webhooks and health check) require authentication using an API key.

**Header:**
```
Authorization: Bearer YOUR_API_KEY
```

**Creating an API Key:**

API keys are generated via the backend and stored in the database. Each key is hashed using bcrypt for security.

---

## Messages

### Send Message

Send a text, media, or template message via WhatsApp.

**Endpoint:** `POST /api/v1/messages`

**Request Body:**

```json
{
  "to": "+1234567890",
  "message_type": "text",
  "content": "Hello from Vibecoded!"
}
```

**Message Types:**

#### Text Message
```json
{
  "to": "+1234567890",
  "message_type": "text",
  "content": "Your message here"
}
```

#### Media Message (Image)
```json
{
  "to": "+1234567890",
  "message_type": "image",
  "media_url": "https://example.com/image.jpg",
  "content": "Optional caption"
}
```

#### Media Message (Document)
```json
{
  "to": "+1234567890",
  "message_type": "document",
  "media_url": "https://example.com/document.pdf",
  "content": "Document title"
}
```

#### Template Message
```json
{
  "to": "+1234567890",
  "message_type": "template",
  "template_name": "welcome_message",
  "template_language": "en",
  "template_params": ["John", "Doe"]
}
```

**Response:** `201 Created`
```json
{
  "id": "msg_abc123",
  "whatsapp_message_id": "wamid.xxx",
  "from_number": "+1234567890",
  "to_number": "+1234567890",
  "direction": "outbound",
  "message_type": "text",
  "content": "Hello from Vibecoded!",
  "status": "sent",
  "timestamp": "2025-11-21T10:30:00Z",
  "created_at": "2025-11-21T10:30:00Z",
  "updated_at": "2025-11-21T10:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid phone number or message content
- `401 Unauthorized` - Missing or invalid API key
- `500 Internal Server Error` - Failed to send message

---

### Get Message

Retrieve a specific message by ID.

**Endpoint:** `GET /api/v1/messages/:id`

**Response:** `200 OK`
```json
{
  "id": "msg_abc123",
  "whatsapp_message_id": "wamid.xxx",
  "from_number": "+1234567890",
  "to_number": "+1234567890",
  "direction": "outbound",
  "message_type": "text",
  "content": "Hello from Vibecoded!",
  "status": "delivered",
  "timestamp": "2025-11-21T10:30:00Z",
  "created_at": "2025-11-21T10:30:00Z",
  "updated_at": "2025-11-21T10:30:00Z"
}
```

**Error Responses:**
- `404 Not Found` - Message not found

---

### List Messages

Get a paginated list of messages with optional filters.

**Endpoint:** `GET /api/v1/messages`

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `limit` (optional) - Items per page (default: 20, max: 100)
- `phone` (optional) - Filter by phone number
- `status` (optional) - Filter by status (sent, delivered, read, failed)
- `direction` (optional) - Filter by direction (inbound, outbound)

**Example:**
```
GET /api/v1/messages?phone=+1234567890&limit=50&page=1
```

**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "msg_abc123",
      "whatsapp_message_id": "wamid.xxx",
      "from_number": "+1234567890",
      "to_number": "+1234567890",
      "direction": "outbound",
      "message_type": "text",
      "content": "Hello!",
      "status": "delivered",
      "timestamp": "2025-11-21T10:30:00Z",
      "created_at": "2025-11-21T10:30:00Z",
      "updated_at": "2025-11-21T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 100,
    "total_pages": 2
  }
}
```

---

### Search Messages

Search messages by content.

**Endpoint:** `GET /api/v1/messages/search`

**Query Parameters:**
- `q` (required) - Search query
- `page` (optional) - Page number
- `limit` (optional) - Items per page

**Example:**
```
GET /api/v1/messages/search?q=order&limit=20
```

**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "msg_abc123",
      "content": "Your order #12345 has been shipped",
      "timestamp": "2025-11-21T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 5,
    "total_pages": 1
  }
}
```

---

## Contacts

### List Contacts

Get a paginated list of contacts.

**Endpoint:** `GET /api/v1/contacts`

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `limit` (optional) - Items per page (default: 20)

**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "cnt_abc123",
      "phone_number": "+1234567890",
      "name": "John Doe",
      "last_message_at": "2025-11-21T10:30:00Z",
      "message_count": 42,
      "unread_count": 3,
      "created_at": "2025-11-20T08:00:00Z",
      "updated_at": "2025-11-21T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

---

### Get Contact

Retrieve a specific contact by ID.

**Endpoint:** `GET /api/v1/contacts/:id`

**Response:** `200 OK`
```json
{
  "id": "cnt_abc123",
  "phone_number": "+1234567890",
  "name": "John Doe",
  "last_message_at": "2025-11-21T10:30:00Z",
  "message_count": 42,
  "unread_count": 3,
  "created_at": "2025-11-20T08:00:00Z",
  "updated_at": "2025-11-21T10:30:00Z"
}
```

**Error Responses:**
- `404 Not Found` - Contact not found

---

### Update Contact

Update contact information.

**Endpoint:** `PATCH /api/v1/contacts/:id`

**Request Body:**
```json
{
  "name": "John Smith"
}
```

**Response:** `200 OK`
```json
{
  "id": "cnt_abc123",
  "phone_number": "+1234567890",
  "name": "John Smith",
  "last_message_at": "2025-11-21T10:30:00Z",
  "message_count": 42,
  "unread_count": 3,
  "created_at": "2025-11-20T08:00:00Z",
  "updated_at": "2025-11-21T11:00:00Z"
}
```

---

### Search Contacts

Search contacts by name or phone number.

**Endpoint:** `GET /api/v1/contacts/search`

**Query Parameters:**
- `q` (required) - Search query
- `page` (optional) - Page number
- `limit` (optional) - Items per page

**Example:**
```
GET /api/v1/contacts/search?q=john
```

**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "cnt_abc123",
      "phone_number": "+1234567890",
      "name": "John Doe",
      "message_count": 42,
      "unread_count": 3
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 2,
    "total_pages": 1
  }
}
```

---

## Templates

### List Templates

Get a paginated list of message templates.

**Endpoint:** `GET /api/v1/templates`

**Query Parameters:**
- `page` (optional) - Page number
- `limit` (optional) - Items per page

**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": "tpl_abc123",
      "name": "welcome_message",
      "language": "en",
      "category": "MARKETING",
      "content": "Welcome {{1}}! Thanks for joining {{2}}.",
      "status": "approved",
      "created_at": "2025-11-20T08:00:00Z",
      "updated_at": "2025-11-20T08:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 10,
    "total_pages": 1
  }
}
```

---

### Get Template

Retrieve a specific template by ID.

**Endpoint:** `GET /api/v1/templates/:id`

**Response:** `200 OK`
```json
{
  "id": "tpl_abc123",
  "name": "welcome_message",
  "language": "en",
  "category": "MARKETING",
  "content": "Welcome {{1}}! Thanks for joining {{2}}.",
  "status": "approved",
  "created_at": "2025-11-20T08:00:00Z",
  "updated_at": "2025-11-20T08:00:00Z"
}
```

---

### Create Template

Create a new message template.

**Endpoint:** `POST /api/v1/templates`

**Request Body:**
```json
{
  "name": "order_confirmation",
  "language": "en",
  "category": "TRANSACTIONAL",
  "content": "Your order {{1}} has been confirmed. Total: ${{2}}"
}
```

**Response:** `201 Created`
```json
{
  "id": "tpl_xyz789",
  "name": "order_confirmation",
  "language": "en",
  "category": "TRANSACTIONAL",
  "content": "Your order {{1}} has been confirmed. Total: ${{2}}",
  "status": "pending",
  "created_at": "2025-11-21T11:00:00Z",
  "updated_at": "2025-11-21T11:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid template data

---

### Update Template

Update an existing template.

**Endpoint:** `PATCH /api/v1/templates/:id`

**Request Body:**
```json
{
  "content": "Your order {{1}} has been confirmed! Total: ${{2}}. Track at {{3}}"
}
```

**Response:** `200 OK`
```json
{
  "id": "tpl_xyz789",
  "name": "order_confirmation",
  "language": "en",
  "category": "TRANSACTIONAL",
  "content": "Your order {{1}} has been confirmed! Total: ${{2}}. Track at {{3}}",
  "status": "pending",
  "created_at": "2025-11-21T11:00:00Z",
  "updated_at": "2025-11-21T11:15:00Z"
}
```

---

### Delete Template

Delete a template (soft delete).

**Endpoint:** `DELETE /api/v1/templates/:id`

**Response:** `204 No Content`

---

## Webhooks

### Verify Webhook

WhatsApp webhook verification endpoint (GET request from Meta).

**Endpoint:** `GET /webhooks/whatsapp`

**Query Parameters:**
- `hub.mode` - Should be "subscribe"
- `hub.verify_token` - Your webhook verify token
- `hub.challenge` - Challenge string to echo back

**Response:** `200 OK`
```
{hub.challenge value}
```

---

### Receive Webhook

Receive incoming webhook events from WhatsApp.

**Endpoint:** `POST /webhooks/whatsapp`

**Headers:**
- `X-Hub-Signature-256` - HMAC SHA256 signature for verification

**Request Body:** (varies by event type)

**Incoming Message Example:**
```json
{
  "object": "whatsapp_business_account",
  "entry": [
    {
      "id": "WHATSAPP_BUSINESS_ACCOUNT_ID",
      "changes": [
        {
          "value": {
            "messaging_product": "whatsapp",
            "metadata": {
              "display_phone_number": "15551234567",
              "phone_number_id": "PHONE_NUMBER_ID"
            },
            "contacts": [
              {
                "profile": {
                  "name": "John Doe"
                },
                "wa_id": "1234567890"
              }
            ],
            "messages": [
              {
                "from": "1234567890",
                "id": "wamid.xxx",
                "timestamp": "1234567890",
                "text": {
                  "body": "Hello!"
                },
                "type": "text"
              }
            ]
          },
          "field": "messages"
        }
      ]
    }
  ]
}
```

**Response:** `200 OK`

---

## System

### Health Check

Check if the API is running and database is accessible.

**Endpoint:** `GET /health`

**No authentication required**

**Response:** `200 OK`
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2025-11-21T11:30:00Z"
}
```

**Error Response:** `503 Service Unavailable`
```json
{
  "status": "unhealthy",
  "database": "disconnected",
  "error": "database connection failed"
}
```

---

## Error Responses

All error responses follow this format:

```json
{
  "code": "ERROR_CODE",
  "message": "Human readable error message",
  "details": {
    "field": "Additional context"
  }
}
```

**Common Error Codes:**
- `INVALID_PHONE_NUMBER` - Phone number format is invalid
- `INVALID_REQUEST` - Request body validation failed
- `UNAUTHORIZED` - Missing or invalid API key
- `NOT_FOUND` - Resource not found
- `DATABASE_ERROR` - Database operation failed
- `WHATSAPP_API_ERROR` - WhatsApp API returned an error

---

## Rate Limiting

- **Default Limit:** 1000 requests per minute per API key
- **Header:** `X-RateLimit-Remaining` shows remaining requests
- **Response on Limit:** `429 Too Many Requests`

---

## Pagination

All list endpoints support pagination with these query parameters:
- `page` - Page number (starts at 1)
- `limit` - Items per page (max 100)

Response includes pagination metadata:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

---

## Phone Number Format

All phone numbers must be in E.164 format:
- Must start with `+`
- Country code required
- No spaces or special characters

**Examples:**
- ✅ `+12345678900`
- ✅ `+919876543210`
- ❌ `12345678900` (missing +)
- ❌ `+1 234-567-8900` (contains spaces/dashes)

---

## Message Status Flow

Messages go through these statuses:
1. `sent` - Message sent to WhatsApp
2. `delivered` - Message delivered to recipient
3. `read` - Message read by recipient
4. `failed` - Message delivery failed

---

For more details, see the [API Design](./API_DESIGN.md) document.
