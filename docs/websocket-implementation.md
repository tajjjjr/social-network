# WebSocket Implementation Documentation

## Overview

This document provides comprehensive documentation for the WebSocket implementation in the social network backend. The WebSocket system enables real-time messaging, notifications, and live updates throughout the application.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Package Structure](#package-structure)
3. [Core Components](#core-components)
4. [Message Types](#message-types)
5. [API Endpoints](#api-endpoints)
6. [Integration Guide](#integration-guide)
7. [Testing](#testing)
8. [Extending the System](#extending-the-system)

## Architecture Overview

The WebSocket implementation follows a clean architecture pattern with dependency injection, making it highly testable and extensible.

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   HTTP Router   │────│  WebSocket       │────│   Database      │
│                 │    │  Manager         │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                    ┌─────────┼─────────┐
                    │         │         │
            ┌───────▼───┐ ┌───▼───┐ ┌───▼────────┐
            │ Session   │ │ Group │ │ Message    │
            │ Resolver  │ │ Query │ │ Persister  │
            └───────────┘ └───────┘ └────────────┘
```

### Key Design Principles

- **Interface-based Design**: Core functionality defined through interfaces for easy testing and mocking
- **Dependency Injection**: All dependencies injected at construction time
- **Thread Safety**: Concurrent connection management with proper synchronization
- **Graceful Error Handling**: Database failures don't interrupt real-time message flow
- **Scalable Architecture**: Designed to handle multiple concurrent connections

## Package Structure

```
backend/internal/websocket/
├── ws.go                    # Core WebSocket manager and client implementation
├── session_resolver.go     # Session-based authentication
├── groups.go              # Group membership management
├── message_persister.go   # Message storage and retrieval
├── notification.go        # Real-time notification system
├── chat_messages.go       # HTTP API handlers for chat functionality
├── integration_test.go    # Comprehensive integration tests
├── ws_test.go            # Unit tests for core functionality
└── websocket_routes_test.go # Route and endpoint testing
```

## Core Components

### 1. WebSocket Manager

The `Manager` is the central component that handles all WebSocket connections and message routing.

```go
type Manager struct {
    clients    map[int64]*Client
    mu         sync.RWMutex
    resolver   SessionResolver
    groupQuery GroupMemberFetcher
    persister  MessagePersister
}
```

**Key Methods:**
- `Register(c *Client)`: Add new client connection
- `Unregister(id int64)`: Remove client connection
- `SendToUser(id int64, msg []byte)`: Send message to specific user
- `BroadcastToGroup(sender int64, groupID int64, msg []byte)`: Broadcast to group members
- `BroadcastToAll(msg []byte)`: Send to all connected users
- `IsOnline(userID int64) bool`: Check user online status

### 2. Client Management

Each WebSocket connection is represented by a `Client` struct:

```go
type Client struct {
    ID        int64
    Conn      *websocket.Conn
    Send      chan []byte
    Connected time.Time
}
```

**Features:**
- 256-buffer channel for non-blocking message sending
- Connection timestamp tracking
- Automatic cleanup on disconnect

### 3. Core Interfaces

The system is built around three main interfaces:

#### SessionResolver
```go
type SessionResolver interface {
    GetUserFromRequest(r *http.Request) (int64, error)
}
```
Handles user authentication via session cookies.

#### GroupMemberFetcher
```go
type GroupMemberFetcher interface {
    GetGroupMemberIDs(groupID int64) ([]int64, error)
}
```
Retrieves group membership for message broadcasting.

#### MessagePersister
```go
type MessagePersister interface {
    SaveMessage(senderID int64, msg *Message) error
}
```
Handles message storage and retrieval with pagination support.

## Message Types

The WebSocket system supports four distinct message types:

### 1. Private Messages (`"private"`)
Direct user-to-user communication.

```json
{
  "type": "private",
  "to": 123,
  "content": "Hello there!",
  "timestamp": 1640995200
}
```

**Behavior:**
- Sent only to the specified recipient
- Persisted to database with sender/receiver IDs
- Requires `to` field with recipient user ID

### 2. Group Messages (`"group"`)
Messages broadcast to all group members.

```json
{
  "type": "group",
  "group_id": "456",
  "content": "Hello everyone!",
  "timestamp": 1640995200
}
```

**Behavior:**
- Broadcast to all accepted group members (excluding sender)
- Persisted to database with group ID
- Requires `group_id` field

### 3. Broadcast Messages (`"broadcast"`)
System-wide announcements to all connected users.

```json
{
  "type": "broadcast",
  "content": "System maintenance in 10 minutes",
  "timestamp": 1640995200
}
```

**Behavior:**
- Sent to all currently connected users
- Not persisted to database (ephemeral)
- Typically used for system announcements

### 4. Notification Messages (`"notification"`)
Server-initiated push notifications.

```json
{
  "type": "notification",
  "subtype": "group_invite",
  "message": "You were invited to join group 'Developers'",
  "group_id": 789,
  "timestamp": 1640995200
}
```

**Behavior:**
- Sent from server to specific users
- Used for real-time alerts (group invites, friend requests, etc.)
- Flexible payload structure for different notification types

## API Endpoints

### WebSocket Endpoint

**`GET /ws`**
- Upgrades HTTP connection to WebSocket
- Requires session-based authentication
- Protected by authentication middleware

### HTTP API Endpoints

#### Chat History
- **`GET /api/messages/private?user={userId}&limit={limit}&offset={offset}`**
  - Retrieve paginated private message history
  - Requires authentication

- **`GET /api/messages/group?group={groupId}&limit={limit}&offset={offset}`**
  - Retrieve paginated group message history
  - Requires authentication

#### Group Management
- **`POST /api/groups/invite`**
  - Send group invitations with real-time notifications
  - Body: `{"group_id": 1, "user_id": 2, "group_name": "Developers"}`

#### Notifications
- **`GET /api/notifications?limit={limit}&offset={offset}`**
  - Retrieve user notification history with read status

- **`POST /api/notifications/read`**
  - Mark all notifications as read for the authenticated user

All endpoints require session-based authentication and return appropriate HTTP status codes.

## Related Documentation

- **[Integration Guide](websocket-integration-guide.md)**: How to integrate WebSocket features into other parts of the application
- **[API Reference](websocket-api-reference.md)**: Complete API documentation for WebSocket endpoints and message formats
- **[Testing Guide](websocket-testing-guide.md)**: Comprehensive testing strategies and examples
- **[Extension Guide](websocket-extension-guide.md)**: How to extend the WebSocket system with new features

## Quick Links

- [Message Types](#message-types) - Understanding different message formats
- [Integration Examples](websocket-integration-guide.md#feature-specific-integration-examples) - Real-world integration patterns
- [Testing Your Integration](websocket-testing-guide.md#testing-your-integration) - How to test WebSocket features
- [Adding New Message Types](websocket-extension-guide.md#adding-new-message-types) - Extending the system
