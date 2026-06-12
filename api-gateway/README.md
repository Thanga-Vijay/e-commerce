# API Gateway Service

**Port:** 8080  
**Language:** Golang  
**Framework:** Gin

## Responsibilities

- **Request Routing** - Route incoming requests to appropriate microservices
- **Authentication** - JWT token validation
- **Rate Limiting** - Protect services from overload
- **Request Logging** - Log all incoming requests
- **Service Discovery** - Dynamic service discovery
- **Load Balancing** - Distribute requests across service instances
- **CORS Handling** - Cross-origin resource sharing

## Architecture

The API Gateway acts as the single entry point for all client requests. It validates JWT tokens, applies rate limiting, and routes requests to the appropriate backend service.

## API Endpoints

The gateway proxies all service endpoints:

- `/api/v1/auth/**` → Auth Service (8081)
- `/api/v1/products/**` → Product Service (8082)
- `/api/v1/cart/**` → Cart Service (8083)
- `/api/v1/wishlist/**` → Wishlist Service (8084)
- `/api/v1/orders/**` → Order Service (8085)
- `/api/v1/payments/**` → Payment Service (8086)
- `/api/v1/inventory/**` → Inventory Service (8087)
- `/api/v1/reports/**` → Reporting Service (8089)

## Configuration

Environment variables:
- `PORT` - Service port (default: 8080)
- `JWT_SECRET` - JWT signing secret
- `RATE_LIMIT` - Requests per minute per IP
- `AUTH_SERVICE_URL` - Auth service URL
- `PRODUCT_SERVICE_URL` - Product service URL
- ... (other service URLs)

## Directory Structure

```
api-gateway/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── middleware/
│   ├── router/
│   └── proxy/
├── pkg/
├── Dockerfile
├── go.mod
└── go.sum
```

## Running Locally

```bash
cd api-gateway
go mod download
go run cmd/main.go
```

## Running with Docker

```bash
docker build -t api-gateway:latest .
docker run -p 8080:8080 api-gateway:latest
```
