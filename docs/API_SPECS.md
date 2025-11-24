# API Specifications

Base URL: http://localhost:8080/api

## Authentication Endpoints

### POST /register

Register a new user account.

Request Body:
```json
{
  "username": "string (required)",
  "email": "string (required, valid email format)",
  "password": "string (required, minimum 6 characters)"
}
```

Success Response (201 Created):
```json
{
  "message": "User registered successfully",
  "token": "string (JWT token)",
  "user": {
    "id": "number",
    "username": "string",
    "email": "string",
    "avatar": "string"
  }
}
```

Error Responses:
- 400 Bad Request: Invalid input or validation error
- 409 Conflict: Username or email already exists
- 500 Internal Server Error: Server error

### POST /login

Authenticate an existing user.

Request Body:
```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

Success Response (200 OK):
```json
{
  "token": "string (JWT token)",
  "user": {
    "id": "number",
    "username": "string",
    "email": "string",
    "avatar": "string"
  }
}
```

Error Responses:
- 400 Bad Request: Invalid input
- 401 Unauthorized: Invalid credentials
- 500 Internal Server Error: Server error

---

## User Endpoints

### GET /users

Get all users. (Public)

Success Response (200 OK):
```json
{
  "users": [
    {
      "id": "number",
      "username": "string",
      "email": "string",
      "avatar": "string",
      "created_at": "string (ISO 8601 datetime)",
      "updated_at": "string (ISO 8601 datetime)"
    }
  ]
}
```

### GET /users/:id

Get a specific user by ID. (Public)

Success Response (200 OK):
```json
{
  "user": {
    "id": "number",
    "username": "string",
    "email": "string",
    "avatar": "string",
    "created_at": "string (ISO 8601 datetime)",
    "updated_at": "string (ISO 8601 datetime)"
  }
}
```

Error Responses:
- 400 Bad Request: Invalid user ID
- 404 Not Found: User not found

### POST /users/:id/block

Block a user. (Protected)

Success Response (200 OK):
```json
{
  "message": "User blocked"
}
```

Error Responses:
- 400 Bad Request: Cannot block yourself
- 401 Unauthorized: Not authenticated
- 409 Conflict: User already blocked

### DELETE /users/:id/block

Unblock a user. (Protected)

Success Response (200 OK):
```json
{
  "message": "User unblocked"
}
```

### GET /blocks

Get list of blocked users. (Protected)

Success Response (200 OK):
```json
{
  "blocked_users": [
    {
      "id": "number",
      "blocker_id": "number",
      "blocked_id": "number",
      "blocked": {
        "id": "number",
        "username": "string"
      },
      "created_at": "string (ISO 8601 datetime)"
    }
  ]
}
```

---

## Room Endpoints

### GET /rooms

Get all chat rooms. (Public)

Success Response (200 OK):
```json
{
  "rooms": [
    {
      "id": "number",
      "name": "string",
      "type": "string (public|private|direct)",
      "created_by": "number",
      "creator": { "id": "number", "username": "string" },
      "members": [],
      "created_at": "string (ISO 8601 datetime)"
    }
  ]
}
```

### GET /rooms/:id

Get a specific room. (Public)

Success Response (200 OK):
```json
{
  "room": {
    "id": "number",
    "name": "string",
    "type": "string",
    "created_by": "number",
    "creator": { "id": "number", "username": "string" },
    "members": []
  }
}
```

### POST /rooms

Create a new room. (Protected)

Request Body:
```json
{
  "name": "string (required)",
  "type": "string (optional, default: public)"
}
```

Success Response (201 Created):
```json
{
  "room": {
    "id": "number",
    "name": "string",
    "type": "string",
    "created_by": "number"
  }
}
```

### POST /rooms/:id/join

Join a room. (Protected)

Success Response (200 OK):
```json
{
  "message": "Joined room successfully",
  "room": { ... }
}
```

Error Responses:
- 404 Not Found: Room not found
- 409 Conflict: Already a member

### POST /rooms/:id/leave

Leave a room. (Protected)

Success Response (200 OK):
```json
{
  "message": "Left room successfully"
}
```

---

## Message Endpoints

### GET /messages

Get all messages with pagination. (Public)

Query Parameters:
- `limit`: Number of messages (default: 50)
- `offset`: Number to skip (default: 0)

Success Response (200 OK):
```json
{
  "messages": [
    {
      "id": "number",
      "user_id": "number",
      "user": { "id": "number", "username": "string" },
      "room_id": "number",
      "room": { "id": "number", "name": "string" },
      "content": "string",
      "edited": "boolean",
      "deleted": "boolean",
      "created_at": "string (ISO 8601 datetime)"
    }
  ]
}
```

### GET /messages/search

Search messages by content. (Public)

Query Parameters:
- `q`: Search query (required)
- `room_id`: Filter by room (optional)
- `limit`: Number of results (default: 50)
- `offset`: Number to skip (default: 0)

Success Response (200 OK):
```json
{
  "messages": [...],
  "query": "string",
  "count": "number"
}
```

### POST /messages

Create a new message. (Protected)

Request Body:
```json
{
  "content": "string (required)",
  "room_id": "number (required)"
}
```

Success Response (201 Created):
```json
{
  "message": {
    "id": "number",
    "user_id": "number",
    "room_id": "number",
    "content": "string",
    "created_at": "string (ISO 8601 datetime)"
  }
}
```

---

## Reaction Endpoints

### GET /messages/:id/reactions

Get reactions for a message. (Public)

Success Response (200 OK):
```json
{
  "reactions": [
    {
      "id": "number",
      "message_id": "number",
      "user_id": "number",
      "user": { "id": "number", "username": "string" },
      "emoji": "string",
      "created_at": "string (ISO 8601 datetime)"
    }
  ],
  "counts": {
    "👍": 2,
    "❤️": 1
  }
}
```

### POST /messages/:id/reactions

Toggle a reaction (add/remove). (Protected)

Request Body:
```json
{
  "emoji": "string (required)"
}
```

Success Response (200 OK):
```json
{
  "message": "Reaction added|removed",
  "added": "boolean",
  "emoji": "string"
}
```

---

## Direct Message Endpoints

### GET /conversations

Get all DM conversations. (Protected)

Success Response (200 OK):
```json
{
  "conversations": [
    {
      "id": "number",
      "user1_id": "number",
      "user1": { "id": "number", "username": "string" },
      "user2_id": "number",
      "user2": { "id": "number", "username": "string" },
      "created_at": "string (ISO 8601 datetime)"
    }
  ]
}
```

### POST /conversations

Start a conversation with a user. (Protected)

Request Body:
```json
{
  "user_id": "number (required)"
}
```

Success Response (200 OK):
```json
{
  "conversation": { ... }
}
```

Error Responses:
- 400 Bad Request: Cannot start conversation with yourself

### GET /conversations/:id/messages

Get messages in a conversation. (Protected)

Query Parameters:
- `limit`: Number of messages (default: 50)
- `offset`: Number to skip (default: 0)

Success Response (200 OK):
```json
{
  "messages": [
    {
      "id": "number",
      "conversation_id": "number",
      "sender_id": "number",
      "sender": { "id": "number", "username": "string" },
      "content": "string",
      "read": "boolean",
      "created_at": "string (ISO 8601 datetime)"
    }
  ]
}
```

### POST /conversations/:id/messages

Send a direct message. (Protected)

Request Body:
```json
{
  "content": "string (required)"
}
```

Success Response (201 Created):
```json
{
  "message": { ... }
}
```

### GET /conversations/unread

Get unread DM count. (Protected)

Success Response (200 OK):
```json
{
  "unread_count": "number"
}
```

---

## Read Receipt Endpoints

### GET /rooms/:room_id/receipts

Get read receipts for a room. (Public)

Success Response (200 OK):
```json
{
  "receipts": [
    {
      "id": "number",
      "user_id": "number",
      "user": { "id": "number", "username": "string" },
      "room_id": "number",
      "last_message_id": "number",
      "last_read_at": "string (ISO 8601 datetime)"
    }
  ]
}
```

### POST /receipts

Mark messages as read. (Protected)

Request Body:
```json
{
  "room_id": "number (required)",
  "message_id": "number (required)"
}
```

Success Response (200 OK):
```json
{
  "message": "Marked as read"
}
```

---

## File Upload Endpoints

### POST /upload

Upload a file. (Protected)

Request: `multipart/form-data`
- `file`: File to upload (required)

Allowed types: jpg, jpeg, png, gif, webp, pdf, txt
Max size: 10MB

Success Response (201 Created):
```json
{
  "message": "File uploaded successfully",
  "url": "/uploads/filename.jpg",
  "filename": "string",
  "size": "number (bytes)"
}
```

Error Responses:
- 400 Bad Request: File type not allowed
- 413 Payload Too Large: File exceeds 10MB

### GET /uploads/:filename

Get an uploaded file. (Public)

Returns the file directly.

---

## WebSocket Endpoint

### GET /ws?token=<JWT>

WebSocket endpoint for real-time chat.

Query Parameters:
- `token`: JWT token (required)

### Message Types (Client → Server)

**Join Room:**
```json
{
  "type": "join_room",
  "room_id": "number"
}
```

**Leave Room:**
```json
{
  "type": "leave_room",
  "room_id": "number"
}
```

**Chat Message:**
```json
{
  "type": "chat",
  "content": "string",
  "room_id": "number (optional, omit for global)"
}
```

**Typing Indicator:**
```json
{
  "type": "typing",
  "room_id": "number (optional)"
}
```

**Stop Typing:**
```json
{
  "type": "stop_typing",
  "room_id": "number (optional)"
}
```

### Message Types (Server → Client)

**Chat Message:**
```json
{
  "type": "chat",
  "content": "string",
  "user_id": "number",
  "username": "string",
  "room_id": "number (if room message)",
  "timestamp": "string (ISO 8601)"
}
```

**User Joined (Global):**
```json
{
  "type": "user_joined",
  "user_id": "number",
  "username": "string"
}
```

**User Left (Global):**
```json
{
  "type": "user_left",
  "user_id": "number",
  "username": "string"
}
```

**Room User Joined:**
```json
{
  "type": "room_user_joined",
  "room_id": "number",
  "user_id": "number",
  "username": "string"
}
```

**Room User Left:**
```json
{
  "type": "room_user_left",
  "room_id": "number",
  "user_id": "number",
  "username": "string"
}
```

**Typing Indicator:**
```json
{
  "type": "typing",
  "user_id": "number",
  "username": "string",
  "room_id": "number (if in room)"
}
```

---

## Health Check

### GET /health

Check API health status. (Public)

Success Response (200 OK):
```json
{
  "status": "ok"
}
```

---

## Authentication

Protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

Token is obtained from /register or /login endpoints.
Token expires after 24 hours.

## Error Response Format

All error responses follow this format:
```json
{
  "error": "string (error description)"
}
```

## Notes

- All datetime values are in ISO 8601 format
- Passwords are hashed using bcrypt before storage
- JWT tokens are signed with HS256 algorithm
- Content-Type for all requests should be application/json (except file uploads)
- CORS is enabled for all origins (development only)
- WebSocket authentication uses query parameter token instead of header
