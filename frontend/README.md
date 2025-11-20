# Vibecoded WA Client - Frontend

🍓 Strawberry-themed WhatsApp Business API client built with React, Vite, and Tailwind CSS.

## Features

- 📱 Real-time messaging interface
- 👥 Contact management
- 📝 Template management
- 🎨 Beautiful Strawberry design system
- ⚡ Fast and responsive
- 🔒 Secure API communication

## Tech Stack

- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool
- **Tailwind CSS** - Styling
- **React Router** - Navigation
- **Axios** - API client
- **Lucide React** - Icons

## Getting Started

### Prerequisites

- Node.js 18+ and npm/yarn
- Backend API running on http://localhost:8080

### Installation

1. Install dependencies:
```bash
npm install
```

2. Create `.env` file:
```bash
cp .env.example .env
```

3. Update `.env` with your API configuration:
```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_API_KEY=your-api-key-here
```

### Development

Start the development server:
```bash
npm run dev
```

The app will be available at http://localhost:5173

### Build

Build for production:
```bash
npm run build
```

Preview production build:
```bash
npm run preview
```

## Project Structure

```
src/
├── components/      # Reusable UI components
├── pages/          # Page components
├── layouts/        # Layout components
├── services/       # API client
├── hooks/          # Custom React hooks
├── lib/            # Utility functions
├── types/          # TypeScript types
└── App.tsx         # Main app component
```

## Design System

This project uses the **Strawberry Theme** design system:

### Colors
- **Primary**: Strawberry Red (#f43f5e)
- **Success**: Leaf Green (#22c55e)
- **Neutral**: Gray shades

### Components
All components follow the design system specified in `docs/design/UI_UX_GUIDE.md`

## API Integration

The frontend communicates with the backend API using Axios. All API calls are centralized in `src/services/api.ts`.

### Authentication
API requests include a Bearer token from the environment variable:
```typescript
Authorization: Bearer ${VITE_API_KEY}
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines.

## License

See [LICENSE](../LICENSE) for license information.

---

**Vibecoded** 🎵
