# Auth Service

**Port:** 8081  
**Language:** Golang  
**Framework:** Gin  
**Database:** PostgreSQL  
**ORM:** GORM

## Responsibilities

- **User Registration** - Create new user accounts
- **User Login** - Authenticate users
- **JWT Generation** - Generate access and refresh tokens
- **Token Refresh** - Refresh expired access tokens
- **Password Reset** - Handle password reset flow
- **Email Verification** - Verify user email addresses
- **RBAC** - Role-based access control
- **User Management** - Update user profiles

## Database Schema

### users
- id (UUID, PK)
- email (unique)
- password_hash
- first_name
- last_name
- role (enum: customer, admin)
- is_verified
- created_at
- updated_at

### refresh_tokens
- id (UUID, PK)
- user_id (FK)
- token (unique)
- expires_at
- created_at

### password_resets
- id (UUID, PK)
- user_id (FK)
- token
- expires_at
- used
- created_at

## API Endpoints

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - Logout user
- `POST /api/v1/auth/forgot-password` - Request password reset
- `POST /api/v1/auth/reset-password` - Reset password
- `POST /api/v1/auth/verify-email` - Verify email
- `GET /api/v1/auth/me` - Get current user
- `PUT /api/v1/auth/profile` - Update profile

## Directory Structure

```
auth-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   ├── service/
│   ├── middleware/
│   └── utils/
├── migrations/
├── pkg/
├── Dockerfile
├── go.mod
└── go.sum
```

## Environment Variables

- `PORT` - Service port (8081)
- `DB_HOST` - PostgreSQL host
- `DB_PORT` - PostgreSQL port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `JWT_SECRET` - JWT signing secret
- `JWT_EXPIRY` - Access token expiry (15m)
- `REFRESH_TOKEN_EXPIRY` - Refresh token expiry (7d)
- `SMTP_HOST` - Email server host
- `SMTP_PORT` - Email server port
- `SMTP_USER` - Email username
- `SMTP_PASSWORD` - Email password

Test