# GoChatApp Frontend

Real-time chat application frontend built with React, Vite, and Tailwind CSS.

## Prerequisites

- Node.js 20.19+ or 22.12+ (required by Vite 7)
- npm or yarn

## Setup

1. Install dependencies:
```bash
npm install
```

2. Configure environment variables:
Copy `.env.example` to `.env` and adjust the values if needed:
```bash
cp .env.example .env
```

Default values:
- `VITE_API_BASE_URL=http://localhost:8080`
- `VITE_WS_URL=ws://localhost:8080/ws`

3. Start the development server:
```bash
npm run dev
```

## Project Structure

```
src/
├── components/     # Reusable UI components
├── pages/          # Page components (Home, Login, Register, Chat)
├── hooks/          # Custom React hooks
├── utils/          # Utility functions (api, websocket)
├── services/       # Service layer for API calls
└── contexts/       # React context providers
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Features

- User authentication (login/register)
- Real-time messaging with WebSocket
- Responsive design with Tailwind CSS
- Protected routes
- Auto-reconnection for WebSocket

## Pages

- `/` - Home page with login/register links
- `/login` - Login page
- `/register` - Registration page
- `/chat` - Chat interface (protected route)

## Note on Node.js Version

This project uses Vite 7, which requires Node.js 20.19+ or 22.12+. If you have an older version of Node.js installed, please upgrade before running the development server.

You can check your Node.js version with:
```bash
node --version
```

To upgrade Node.js, visit: https://nodejs.org/
