# Features Roadmap

Reference: [Whatomate](https://shridarpatil.github.io/whatomate/) | [GitHub](https://github.com/shridarpatil/whatomate)

## Current State

Our project already has basic implementations for:
- Messages (send/receive text, media, templates)
- Contacts (CRUD, search)
- Templates (CRUD)
- Webhooks (Meta verification + inbound events)
- Health endpoint
- API key authentication
- SQLite/PostgreSQL database support
- React frontend (Messages, Contacts, Templates pages)

---

## Phase 1: Real-Time & Chat Experience

**Goal**: Transform from API-only to a live chat platform.

- [ ] **WebSocket server** — Hub-based pattern for pushing real-time updates to connected clients
  - Events: `message:new`, `message:status`, `contact:new`, `contact:updated`
  - JWT-authenticated connections (`ws://host/ws?token=TOKEN`)
- [ ] **Message status tracking** — Full lifecycle: pending → sent → delivered → read → failed
- [ ] **Reply threading** — `is_reply` flag, `reply_to_message_id` on messages
- [ ] **Read receipts** — Mark messages as read via API (`PUT /messages/:id/read`)
- [ ] **All message types** — Extend beyond text/media/template:
  - Interactive messages (buttons up to 3, CTA URLs)
  - Reaction messages
  - Location messages
  - Contact card messages
- [ ] **Chat UI overhaul** — Real-time message feed, typing indicators, message status icons
- [ ] **Contact info panel** — Side panel showing contact details, tags, metadata, notes
- [ ] **Service window tracking** — `last_inbound_at` for 24-hour messaging window awareness

---

## Phase 2: Multi-Tenant & User Management

**Goal**: Support multiple organizations and proper user management.

- [ ] **Organization model** — Create, switch, and manage organizations with isolated data
  - Organization settings: name, timezone, date_format, mask_phone_numbers
  - Slug-based identification
- [ ] **User management** — Full user lifecycle
  - User CRUD (`GET/POST/PUT/DELETE /api/users`)
  - Profile endpoint (`GET /api/me`)
  - Password change with current password verification
  - Availability status tracking with logging
  - Users can belong to multiple organizations
- [ ] **JWT authentication** — Replace or supplement API key auth
  - Access tokens (configurable expiry, default 15 min)
  - Refresh tokens (default 7 days)
  - httpOnly cookie storage
  - Organization-scoped tokens
- [ ] **API key improvements** — `whm_` prefix format, bcrypt hashing, optional expiration, prefix-based identification
- [ ] **Organization switching** — `POST /api/auth/switch-org` to change active context
- [ ] **Organization members** — Add/remove users, assign roles per org

---

## Phase 3: Roles & Permissions (RBAC)

**Goal**: Granular access control across the platform.

- [ ] **Permission system** — `resource:action` format (e.g., `contacts:read`, `campaigns:execute`)
  - Actions: read, write, delete, sync, execute, import, export, pickup, assign
  - 30+ permission resources
- [ ] **System roles** (non-deletable):
  - Admin — All permissions
  - Manager — All except user/role/SSO management
  - Agent — Chat, transfers, own analytics, canned responses read-only
- [ ] **Custom roles** — Organization-specific roles with arbitrary permission sets
- [ ] **Permission enforcement** — Middleware for API-level checks, 403 responses
- [ ] **Frontend enforcement** — Hide UI elements based on permissions, redirect unauthorized access
- [ ] **Roles API** — `GET/POST/PUT/DELETE /api/roles`, `GET /api/permissions`

---

## Phase 4: Teams & Agent Assignment

**Goal**: Structure agents into teams for organized routing.

- [ ] **Team model** — Name, description, assignment strategy, active/inactive
- [ ] **Team roles** — Manager (manages settings, views all transfers) and Agent (handles transfers)
- [ ] **Assignment strategies**:
  - Round robin (sequential)
  - Load balanced (fewest active transfers)
  - Manual (queue for pickup)
- [ ] **Contact assignment** — Assign contacts to specific agents
  - `PUT /api/contacts/:id/assign`
  - Agents see only assigned contacts; admins/managers see all
- [ ] **Teams API** — `GET/POST/PUT/DELETE /api/teams`, member management endpoints

---

## Phase 5: Chatbot Automation

**Goal**: Intelligent automated responses and conversation flows.

### 5a: Core Chatbot Settings
- [ ] **Enable/disable per WhatsApp account**
- [ ] **Greeting message** — With interactive buttons (1-3 = quick reply, 4-10 = list menu)
- [ ] **Fallback message** — For unmatched inputs, with buttons
- [ ] **Session timeout** — Default 30 minutes
- [ ] **Excluded phone numbers** — Skip automation for specific numbers
- [ ] **Business hours** — Per-day schedule with out-of-hours message

### 5b: Keyword Rules
- [ ] **Match types** — exact, contains, starts_with, regex
- [ ] **Case sensitivity toggle**
- [ ] **Priority ordering**
- [ ] **Response types** — text, template, media, flow, script, transfer
- [ ] **Active date range** — `active_from`, `active_until`
- [ ] **Conditions** — Conditional matching rules
- [ ] **Keywords API** — `GET/POST/PUT/DELETE /api/chatbot/keywords`

### 5c: Conversation Flows
- [ ] **Flow builder** — Visual drag-and-drop step editor
- [ ] **Step types** — text, template, script, api_fetch, buttons, transfer, whatsapp_flow
- [ ] **Input types** — none, text, number, email, phone, date, select, button
- [ ] **Input validation** — Regex patterns
- [ ] **Variable storage** — `store_as` for capturing responses
- [ ] **Conditional branching** — Option-to-step mapping
- [ ] **API integration** — Configurable URL, method, headers, body, response_path
- [ ] **Completion actions** — none, webhook, create_record
- [ ] **Flows API** — `GET/POST/PUT/DELETE /api/chatbot/flows`

### 5d: AI Integration
- [ ] **Provider support** — OpenAI (GPT-4o), Anthropic (Claude), Google AI (Gemini)
- [ ] **Configurable** — model, temperature, max tokens, system prompt
- [ ] **Conversation history** — Toggle + history limit
- [ ] **AI Contexts** — Knowledge bases (static text or API-sourced) with trigger keywords
- [ ] **Encrypted API key storage**

---

## Phase 6: Agent Transfers & SLA

**Goal**: Route conversations to human agents with SLA tracking.

- [ ] **Transfer system** — Sources: manual, flow, keyword, chatbot_disabled
  - Statuses: active, resumed, expired
  - Queue management with team-based routing
- [ ] **SLA tracking**:
  - Response time SLA (default 15 min)
  - Resolution time SLA (default 60 min)
  - Escalation time (default 30 min) with notification user IDs
  - Warning message when SLA breached
- [ ] **Auto-close** — Stale transfers after configurable hours (default 24h)
- [ ] **Client inactivity** — Reminder after X minutes, auto-close after Y minutes
- [ ] **Transfers API** — `GET/POST/PUT /api/chatbot/transfers`

---

## Phase 7: Canned Responses & Custom Actions

**Goal**: Productivity tools for agents.

### 7a: Canned Responses
- [ ] **Pre-written snippets** — Name, content, shortcut, category
- [ ] **Categories** — greeting, support, sales, closing, general, custom
- [ ] **Dynamic placeholders** — `{{contact_name}}`, `{{phone_number}}`
- [ ] **Slash command access** — `/shortcut` in chat
- [ ] **Usage tracking** — Auto-increment counter, sort by frequency
- [ ] **API** — `GET/POST/PUT/DELETE /api/canned-responses`

### 7b: Custom Actions
- [ ] **Action types**:
  - Webhook — POST/GET/PUT/PATCH to external APIs with templates
  - URL — Open external URLs with contact parameters
  - JavaScript — Client-side code execution with context
- [ ] **Template variables** — `{{contact.id}}`, `{{contact.phone_number}}`, `{{user.name}}`, etc.
- [ ] **Custom icons** — Lucide icon set
- [ ] **SSRF protection** — Block internal addresses, DNS rebinding protection
- [ ] **API** — `GET/POST/PUT/DELETE /api/custom-actions`, `POST /api/custom-actions/:id/execute`

---

## Phase 8: Campaigns (Bulk Messaging)

**Goal**: Send bulk messages using approved templates.

- [ ] **Campaign lifecycle** — draft → scheduled → queued → processing → paused → completed/cancelled/failed
- [ ] **Recipient management** — Manual entry or CSV upload
  - Smart column name recognition (phone/mobile, name/contact_name, param1/variable1)
  - Duplicate detection
  - 50-row preview
- [ ] **Header media support** — image, video, document
- [ ] **Template variable mapping**
- [ ] **Scheduling** — Specific date/time
- [ ] **Pause/resume/cancel controls**
- [ ] **Rate limiting** — WhatsApp compliance
- [ ] **Delivery tracking** — Total, sent, delivered, read, failed counts per recipient
- [ ] **Retry support** — For failed messages
- [ ] **Campaign API** — Full CRUD + recipients import + start/pause/cancel

---

## Phase 9: WhatsApp Flows

**Goal**: Rich interactive UI experiences native to WhatsApp.

- [ ] **Flow builder** — Visual drag-and-drop, multi-screen journeys
- [ ] **Native UI components** — Text input, dropdown, date picker, radio, checkboxes, action buttons
- [ ] **JSON editor** — Advanced customization
- [ ] **Flow lifecycle** — DRAFT → PUBLISHED → DEPRECATED (+ BLOCKED)
- [ ] **Meta integration** — Save to Meta with validation, sync from Meta
- [ ] **Preview** before deployment
- [ ] **Flows API** — `GET/POST/PUT/DELETE /api/flows` + save-to-meta, publish, deprecate, sync

---

## Phase 10: Templates Enhancement

**Goal**: Full template management with Meta sync.

- [ ] **Two-way Meta sync** — Pull latest templates, push new/updated
- [ ] **All template types** — MARKETING, UTILITY, AUTHENTICATION
- [ ] **Component support** — HEADER (text/image/document/video), BODY, FOOTER, BUTTONS (CTA/quick reply)
- [ ] **Parameter formats** — Positional (`{{1}}`) and named (`{{customer_name}}`)
- [ ] **Status lifecycle** — PENDING → APPROVED / REJECTED
- [ ] **Sample values** — For Meta approval submissions
- [ ] **Publish endpoint** — Submit for Meta approval
- [ ] **Sync endpoint** — `POST /api/templates/sync` returning synced/new/updated counts

---

## Phase 11: Analytics Dashboard

**Goal**: Comprehensive analytics and reporting.

- [ ] **Dashboard stats** — Period-based (today, week, month, year) with account filter
- [ ] **Custom widgets** — Configurable data sources, aggregation, chart types (line, bar, pie)
  - Grid-based positioning
  - Team sharing
  - Color theming
- [ ] **Message analytics** — Volume, delivery rate, read rate, type breakdown
  - Granularity: hour, day, week, month
- [ ] **Agent analytics** — Transfers handled, resolution times, queue metrics, break time
- [ ] **Chatbot analytics** — Conversations, resolution rate, AI usage (tokens, costs)
- [ ] **Campaign performance** — Active, completed, engagement metrics
- [ ] **Data retention** — Hourly (30 days), daily (1 year), monthly (indefinite)
- [ ] **Analytics API** — `/api/analytics/dashboard`, `/api/analytics/messages`, `/api/analytics/chatbot`

---

## Phase 12: SSO & Advanced Auth

**Goal**: Enterprise authentication options.

- [ ] **OAuth2 SSO** — Google, Microsoft, GitHub, Facebook, custom OIDC
- [ ] **Auto-create users** on first SSO login
- [ ] **Allowed email domains** — Restrict SSO to specific domains
- [ ] **Default role assignment** for SSO users
- [ ] **SSO API** — `GET /api/auth/sso/providers`, `GET /api/auth/sso/:provider/init`, callback endpoint
- [ ] **Rate limiting** on auth endpoints

---

## Phase 13: Outbound Webhooks

**Goal**: Notify external systems of events.

- [ ] **Subscribable events**:
  - `message.incoming`, `message.outgoing`, `message.sent`
  - `contact.created`
  - `transfer.created`, `transfer.assigned`, `transfer.resumed`
- [ ] **HMAC signature verification** — Configurable secret
- [ ] **Custom headers** per webhook
- [ ] **Active/inactive toggle**
- [ ] **SSRF protection**
- [ ] **Webhooks API** — `GET/POST/PUT/DELETE /api/webhooks`, `GET /api/webhooks/events`

---

## Phase 14: Voice Calling & IVR

**Goal**: WhatsApp voice calling with IVR system.

- [ ] **Incoming calls** — Webhook-based WhatsApp calling integration
- [ ] **IVR system** — Multi-level menus with audio greetings (file or TTS)
  - DTMF digit routing
  - Action types: transfer, submenu, jump, return, replay, terminate
  - Configurable timeout and max retries
- [ ] **Call transfers** — Hold music, agent notification, first-accept bridging, team routing
- [ ] **Call recording** — OGG/Opus format, S3 upload, duration tracking
- [ ] **Call logs** — Full call history with IVR path traversal
- [ ] **Outgoing calls** — Agent-initiated outbound from chat
- [ ] **WebRTC** — Audio bridging with STUN/TURN support
- [ ] **Text-to-Speech** — Piper integration for offline TTS with caching

---

## Phase 15: WhatsApp Accounts & Business Profile

**Goal**: Multi-account support and business profile management.

- [ ] **Multiple WhatsApp accounts** per organization
- [ ] **Default incoming/outgoing** account designation
- [ ] **Connection status** — pending, active, disconnected, suspended
- [ ] **Quality rating** — GREEN, YELLOW, RED monitoring
- [ ] **Messaging tier** tracking (e.g., TIER_1K)
- [ ] **Auto read receipts** toggle
- [ ] **Access token encryption** at rest
- [ ] **App Secret** for webhook signature verification (HMAC-SHA256)
- [ ] **Test connectivity** endpoint — `POST /api/accounts/:id/test`
- [ ] **Business profile management** — Display name, description, address, etc.

---

## Phase 16: Additional Features

- [ ] **Contact tags** — Color-coded (blue, red, green, yellow, purple, gray), org-scoped
- [ ] **Conversation notes** — Private internal notes on contacts (not visible to customers)
- [ ] **Contact import/export** — Bulk CSV import, export functionality
- [ ] **Product catalog** — WhatsApp catalog integration for commerce
- [ ] **Notification rules** — Automated notifications with triggers (webhook, scheduler, API)
- [ ] **Internationalization (i18n)** — Multi-language frontend support
- [ ] **Contact metadata** — Nested JSON metadata with rich rendering

---

## Phase 17: Infrastructure & Operations

- [ ] **Redis integration** — Queue system, rate limiting, caching
- [ ] **Background worker** — Separate worker process for async tasks (`./server worker -workers=4`)
- [ ] **Rate limiting** — Per-endpoint with Redis fixed-window counters, rate limit headers
- [ ] **CSRF protection** middleware
- [ ] **AES-256 encryption** for secrets at rest (API keys, tokens)
- [ ] **Embedded frontend** — Single binary with embedded static assets
- [ ] **CLI modes** — server, worker, version
- [ ] **Configuration via TOML** — Structured config file support alongside env vars

---

## Priority Recommendation

**High value, build first**: Phases 1-4 (real-time chat, users, RBAC, teams)
**Core automation**: Phases 5-6 (chatbot, transfers)
**Agent productivity**: Phases 7-8 (canned responses, campaigns)
**Advanced features**: Phases 9-14 (flows, analytics, SSO, calling)
**Polish**: Phases 15-17 (multi-account, catalog, i18n, infrastructure)
