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

## User Endpoints

### GET /users

Get all users.

Request Headers:
None required (public endpoint)

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

Error Responses:
- 500 Internal Server Error: Failed to fetch users

### GET /users/:id

Get a specific user by ID.

URL Parameters:
- id: User ID (number)

Request Headers:
None required (public endpoint)

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
- 500 Internal Server Error: Server error

## Message Endpoints

### GET /messages

Get all messages with pagination.

Query Parameters:
- limit: Number of messages to return (default: 50)
- offset: Number of messages to skip (default: 0)

Request Headers:
None required (public endpoint)

Success Response (200 OK):
```json
{
  "messages": [
    {
      "id": "number",
      "user_id": "number",
      "user": {
        "id": "number",
        "username": "string",
        "email": "string",
        "avatar": "string"
      },
      "room_id": "number",
      "room": {
        "id": "number",
        "name": "string",
        "type": "string"
      },
      "content": "string",
      "edited": "boolean",
      "deleted": "boolean",
      "created_at": "string (ISO 8601 datetime)",
      "updated_at": "string (ISO 8601 datetime)"
    }
  ]
}
```

Error Responses:
- 500 Internal Server Error: Failed to fetch messages

### POST /messages

Create a new message.

Request Headers:
None required currently (will require Authorization header in future)

Request Body:
```json
{
  "content": "string (required)",
  "user": "string (temporary field)"
}
```

Success Response (201 Created):
```json
{
  "message": "Message sent"
}
```

Error Responses:
- 400 Bad Request: Invalid input
- 500 Internal Server Error: Failed to create message

## WebSocket Endpoint

### GET /ws

WebSocket endpoint for real-time chat.

Status: Not implemented yet (returns 501)

## Health Check

### GET /health

Check API health status.

Success Response (200 OK):
```json
{
  "status": "ok"
}
```

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
- Content-Type for all requests should be application/json
- CORS is enabled for all origins (development only)
