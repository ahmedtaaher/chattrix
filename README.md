# Chattrix — Real-Time Chat System

A scalable real-time chat backend built with **Golang**, featuring WebSocket communication, message delivery tracking, reactions, mentions, and notifications.

---

## Features

### Authentication

* JWT-based register & login
* Secure protected routes

---

### Messaging System

* Send messages (text, image, file, voice)
* Reply to messages
* Forward messages
* Attachments support

---

### Real-Time Communication

* WebSocket-based messaging
* Typing & stop typing indicators
* Instant message delivery

---

### Message Status

* Sent
* Delivered
* Seen

---

### Reactions

* Add/remove reactions (toggle)
* Aggregated reaction counts per message

---

### Mentions

* Detect `@username` in messages
* Notify mentioned users

---

### Notifications

* In-app notifications (stored in DB)
* Real-time notification delivery
* Mention notifications
* Mute chat support

---

### Chat Management

* Create chats (1:1 / group)
* Add/remove users
* Leave chat
* Mute chats
* Pin chats
* Role management

---

### Search

* Search users
* Search chats
* Search messages (ready for full-text)

---

## Tech Stack

* **Language**: Go (Golang)
* **Framework**: Gin
* **Real-Time**: Gorilla WebSocket
* **Database**: PostgreSQL
* **ORM**: GORM
* **Authentication**: JWT
* **Architecture**: N-Tier (Handler → Service → Repository)
* **Containerization**: Docker 

---

## Project Structure

```text
chattrix/
├── handler/           # HTTP & WebSocket handlers
├── service/           # Business logic
├── repository/        # Database layer
├── models/            # Entities
├── dto/               # Data transfer objects
├── mapper/            # Entity ↔ DTO mapping
├── websocket/         # WebSocket hub & handler
├── middleware/        # Auth middleware
├── utils/             # Helpers (JWT, etc.)
├── db/                # Database connection
├── routes/            # Route definitions
├── main.go            # Entry point
└── docker-compose.yml
```

---

