# GoChatApp - Project Task List

This document outlines all tasks for building the real-time chat application using Go-Gin (backend) and React (frontend).

---

## 🔧 BACKEND TASKS

### Phase 1: Database & Models

- [X] Choose and install database (PostgreSQL/MongoDB/SQLite)
- [X] Create database connection and configuration
- [X] Set up migration system (if using SQL)
- [X] Create User model/schema (username, email, password_hash, avatar, created_at)
- [X] Create Message model/schema (id, user_id, room_id, content, timestamp, edited, deleted)
- [X] Create Room/Channel model/schema (id, name, type, created_by, members)
- [X] Add database seeding for development/testing

### Phase 2: Authentication & Security

- [ ] Install bcrypt library for password hashing
- [ ] Implement JWT token generation (using jwt-go or similar)
- [ ] Create JWT validation middleware
- [ ] Complete `/register` endpoint with password hashing and user creation
- [ ] Complete `/login` endpoint with credential validation and JWT generation
- [ ] Add token refresh endpoint
- [ ] Add authentication middleware for protected routes
- [ ] Implement logout/token invalidation

### Phase 3: REST API Endpoints

- [ ] Update `/messages` GET to return real data from database with pagination
- [ ] Update `/messages` POST to save messages to database
- [ ] Create `/messages/:id` PUT for editing messages
- [ ] Create `/messages/:id` DELETE for deleting messages
- [ ] Update `/users` GET to return real users from database
- [ ] Add `/users/:id` GET for user profile
- [ ] Add `/users/:id` PUT for updating profile
- [ ] Create `/rooms` GET to list all chat rooms
- [ ] Create `/rooms` POST to create new room
- [ ] Create `/rooms/:id/join` POST to join room
- [ ] Create `/rooms/:id/leave` POST to leave room
- [ ] Add `/users/search` GET for user search

### Phase 4: WebSocket Implementation

- [ ] Install gorilla/websocket library
- [ ] Implement WebSocket upgrade in `handleWebSocket` function
- [ ] Create Hub/Manager struct to track active connections
- [ ] Implement client registration (on connect)
- [ ] Implement client unregistration (on disconnect)
- [ ] Create message broadcasting logic (to all clients in room)
- [ ] Add ping/pong heartbeat mechanism
- [ ] Handle different WebSocket message types (chat, typing, etc.)
- [ ] Implement room-based message routing
- [ ] Add online/offline status broadcasting

### Phase 5: Advanced Features

- [ ] Add file upload endpoint (for images/attachments)
- [ ] Implement typing indicator broadcasting via WebSocket
- [ ] Add message reactions endpoint
- [ ] Implement direct messaging (1-on-1 chat)
- [ ] Add message search endpoint
- [ ] Implement user blocking/muting
- [ ] Add read receipts tracking

### Phase 6: Production Readiness

- [ ] Add rate limiting middleware
- [ ] Implement request validation and sanitization
- [ ] Add structured logging (logrus or zap)
- [ ] Set up error handling middleware
- [ ] Add CORS configuration for production
- [ ] Create environment variable configuration
- [ ] Write unit tests for handlers
- [ ] Write integration tests for WebSocket
- [ ] Create Dockerfile
- [ ] Create docker-compose.yml (with database)
- [ ] Write API documentation

---

## ⚛️ FRONTEND TASKS

### Phase 1: Project Setup & Routing

- [X] Initialize React app (Vite/Create React App)
- [X] Install dependencies (react-router, axios, socket.io-client/websocket)
- [X] Set up routing (/, /login, /register, /chat)
- [X] Configure API base URL (environment variables)
- [X] Set up Tailwind CSS or CSS framework
- [X] Create project folder structure (components, pages, hooks, utils)

### Phase 2: Authentication UI

- [ ] Create Login page component
- [ ] Create Registration page component
- [ ] Add form validation (email, password strength)
- [ ] Implement login API call
- [ ] Implement register API call
- [ ] Store JWT token (localStorage or cookies)
- [ ] Create ProtectedRoute component
- [ ] Implement auto-login on page refresh
- [ ] Add logout functionality
- [ ] Create auth context/state management

### Phase 3: Main Chat Interface

- [ ] Create main chat layout (sidebar + message area)
- [ ] Build MessageList component
- [ ] Build Message component (individual message bubble)
- [ ] Build MessageInput component (text input + send button)
- [ ] Build UserSidebar component (online users list)
- [ ] Build RoomList component (available chat rooms)
- [ ] Add responsive design for mobile
- [ ] Create ChatHeader component (room name, users count)

### Phase 4: WebSocket Integration

- [ ] Create WebSocket connection hook/utility
- [ ] Connect to WebSocket on login
- [ ] Handle incoming messages and update UI
- [ ] Send messages via WebSocket
- [ ] Implement auto-reconnection logic
- [ ] Show connection status indicator
- [ ] Handle WebSocket errors gracefully
- [ ] Disconnect WebSocket on logout

### Phase 5: Message Features

- [ ] Display message timestamps (formatted)
- [ ] Implement infinite scroll for message history
- [ ] Add message editing UI
- [ ] Add message deletion UI
- [ ] Show message delivery/read status
- [ ] Differentiate own messages vs others (alignment, color)
- [ ] Add message reactions UI (emoji picker)
- [ ] Implement message search UI

### Phase 6: Real-Time Features

- [ ] Display typing indicators ("User is typing...")
- [ ] Show online/offline status (green/grey dots)
- [ ] Update user list in real-time
- [ ] Add notification system (browser notifications)
- [ ] Add sound notifications for new messages
- [ ] Show new message badge on inactive rooms

### Phase 7: Advanced Features

- [ ] Add emoji picker to message input
- [ ] Implement file/image upload UI
- [ ] Display uploaded images in chat
- [ ] Add user profile modal/page
- [ ] Implement avatar upload
- [ ] Add dark mode toggle
- [ ] Create settings page
- [ ] Add direct messaging UI

### Phase 8: Polish & Testing

- [ ] Add loading states
- [ ] Add error handling and error messages
- [ ] Improve accessibility (ARIA labels, keyboard navigation)
- [ ] Add animations/transitions
- [ ] Test on different browsers
- [ ] Test responsive design on mobile devices
- [ ] Optimize performance (React.memo, lazy loading)
- [ ] Write component tests (Jest/Vitest)

---

## 🤝 COORDINATION POINTS

### API Contract Agreement (Do Together)

- [ ] Define all API endpoints and their request/response formats
- [ ] Document WebSocket event types and payloads
- [ ] Agree on error response format
- [ ] Define authentication header format

### Integration Testing

- [ ] Test login/register flow end-to-end
- [ ] Test real-time messaging works between frontend and backend
- [ ] Test file upload functionality
- [ ] Test WebSocket reconnection behavior

### Deployment

- [ ] Deploy backend to hosting service (Heroku/Railway/Render)
- [ ] Deploy frontend to Vercel/Netlify
- [ ] Configure CORS for production domains
- [ ] Set up CI/CD pipeline

---

## 📝 Notes

- Backend developer focuses on all Backend tasks
- Frontend developer focuses on all Frontend tasks
- Coordinate regularly on API contracts and integration points
- Test integration points frequently to catch issues early
- Update this document as tasks are completed or new tasks are identified

---

## 🎯 Priority Order

1. **Phase 1 (Both)**: Database setup + React project setup
2. **Phase 2 (Both)**: Authentication backend + Authentication UI
3. **Phase 3 & 4 (Backend) + Phase 3 & 4 (Frontend)**: API + WebSocket + Chat UI
4. **Phase 5+**: Advanced features and polish
