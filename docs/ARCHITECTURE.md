# GophKeeper Architecture

## Overview

GophKeeper is a client-server password manager built with Go. The system consists of a server that stores encrypted data and a CLI client that allows users to manage their credentials.

## Architecture Components

### Server (`cmd/server`)

The server provides a REST API for authentication and data management:

- **Authentication**: User registration and login with JWT tokens
- **Data Storage**: PostgreSQL database for persistent storage
- **API Endpoints**:
  - `POST /api/register` - Register new user
  - `POST /api/login` - Login user
  - `GET /api/entries` - List all entries
  - `POST /api/entries` - Create new entry
  - `GET /api/entries/{id}` - Get entry by ID
  - `PUT /api/entries/{id}` - Update entry
  - `DELETE /api/entries/{id}` - Delete entry
  - `GET /api/sync` - Sync entries

### Client (`cmd/client`)

CLI application for managing passwords:

- User authentication
- Create/Read/Update/Delete operations on entries
- Synchronization with server
- Local caching of credentials

### Storage Layer (`internal/storage`)

Database abstraction layer:

- User management (create, authenticate, retrieve)
- Entry CRUD operations
- Data synchronization support

### Crypto (`internal/crypto`)

Security utilities:

- AES-256-GCM encryption for data
- PBKDF2 for key derivation
- bcrypt for password hashing

## Security

### Data Encryption

All sensitive data is encrypted before storage:

1. User provides a master password
2. Data is encrypted with AES-256-GCM
3. Encryption key is derived using PBKDF2 with 4096 iterations
4. Encrypted data is base64 encoded for storage

### Authentication

- Passwords are hashed using bcrypt
- JWT tokens are issued for authenticated sessions
- Tokens include user ID and expiration time
- All protected endpoints require valid JWT

### Data Isolation

- Each user's data is isolated by user_id
- Database foreign keys ensure data integrity
- All queries filter by user_id

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Entries Table

```sql
CREATE TABLE entries (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    data TEXT NOT NULL,
    metadata TEXT,
    version INTEGER DEFAULT 1,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Data Types

### Login Entry

```json
{
  "site": "example.com",
  "login": "user@example.com",
  "password": "secret",
  "notes": "Personal account"
}
```

### Text Entry

```json
{
  "content": "Any text data",
  "notes": "Optional notes"
}
```

### Binary Entry

```json
{
  "filename": "document.pdf",
  "content": "base64-encoded-binary-data",
  "notes": "Important document"
}
```

### Card Entry

```json
{
  "number": "4111111111111111",
  "holder": "John Doe",
  "expiry_month": 12,
  "expiry_year": 2025,
  "cvv": "123",
  "notes": "Primary card"
}
```

## Deployment

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Docker (optional)

### Local Development

1. Start PostgreSQL:
   ```bash
   docker-compose up -d
   ```

2. Run server:
   ```bash
   make run-server
   ```

3. Run client:
   ```bash
   make run-client
   ```

### Production

1. Build binaries:
   ```bash
   make build
   ```

2. Configure environment variables
3. Run server and client binaries

## Testing

Run all tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```
