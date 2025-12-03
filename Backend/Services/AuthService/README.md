# Authentication Service

A production-ready authentication microservice for the AREA project with user registration, login, and JWT-based authentication.

## 🏗️ Architecture

This authentication service provides:

- **Go 1.22**: High-performance backend with native concurrency
- **PostgreSQL 16**: Robust relational database with ACID guarantees
- **JWT Authentication**: Secure token-based authentication with 24-hour expiry
- **bcrypt**: Password hashing with industry-standard security
- **Docker & Docker Compose**: Containerized deployment
- **OpenAPI 3.0**: Complete API specification and documentation
- **RESTful API**: Clean HTTP endpoints with JSON responses

## 📁 Project Structure

```
.
├── app/                      # Application source code
│   ├── cmd/
│   │   └── main.go          # Main application entry point
│   ├── internal/
│   │   ├── auth/            # Authentication utilities
│   │   │   ├── jwt.go      # JWT token generation/validation
│   │   │   └── password.go  # Password hashing/checking
│   │   ├── config/          # Configuration management
│   │   │   └── config.go
│   │   ├── db/              # Database connection
│   │   │   └── db.go
│   │   ├── domain/          # Domain models
│   │   │   └── user.go
│   │   ├── http/            # HTTP handlers and routing
│   │   │   └── router.go
│   │   ├── repository/      # Data access layer
│   │   │   └── user_postgres.go
│   │   └── service/         # Business logic
│   │       └── auth_service.go
│   ├── go.mod              # Go module dependencies
│   └── go.sum              # Go dependency checksums
├── db/                      # Database configuration
│   └── init/               # Database initialization SQL files
│       └── 01_create_tables.sql
├── docker-compose.yml       # Docker Compose configuration
├── Dockerfile               # Multi-stage Docker build
├── Makefile                 # Common commands
├── openapi.yaml             # OpenAPI 3.0 specification
├── .env                     # Environment variables (do not commit)
├── .gitignore               # Git ignore patterns
└── README.md                # This file
```

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.22+ (for local development)
- Make (optional, for convenience commands)

### 1. Initialize Environment

```bash
# Copy environment variables
cp .env.example .env

# Or use make
make init
```

### 2. Start with Docker

```bash
# Build and start containers
docker-compose up -d

# Or use make
make docker-up
```

The API will be available at `http://localhost:8080`

### 3. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Register a new user
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","username":"johndoe","password":"securepass123"}'

# Response includes user data and JWT token:
# {
#   "success": true,
#   "data": {
#     "user": {
#       "id": 1,
#       "email": "john@example.com",
#       "username": "johndoe",
#       "created_at": "2025-01-15T10:30:00Z",
#       "updated_at": "2025-01-15T10:30:00Z"
#     },
#     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
#   }
# }

# Login with email or username
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"emailOrUsername":"johndoe","password":"securepass123"}'

# Get current user profile (requires authentication)
curl http://localhost:8080/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

## 🛠️ Development

### Local Development (without Docker)

```bash
# Navigate to app directory
cd app/

# Install dependencies
go mod download

# Run PostgreSQL (required)
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=microservice_db \
  -p 5432:5432 \
  postgres:16-alpine

# Run the application
go run main.go

# Or use make from root directory
cd ..
make run
```

### Available Make Commands

```bash
make help           # Show all available commands
make build          # Build the Go binary
make run            # Run locally
make test           # Run tests
make clean          # Remove build artifacts
make docker-build   # Build Docker image
make docker-up      # Start Docker containers
make docker-down    # Stop Docker containers
make docker-logs    # View container logs
make docker-restart # Restart containers
make deps           # Download dependencies
make fmt            # Format code
make lint           # Run linter
```

## 📡 API Endpoints

### Health Check
- **GET** `/health` - Check service health

### Authentication
- **POST** `/auth/register` - Register a new user
  - **Body**: `{ "email": string, "username": string, "password": string }`
  - **Returns**: User object + JWT token
  - **Validations**:
    - Email: Valid email format, unique
    - Username: 3-20 alphanumeric characters (including underscore), unique
    - Password: Minimum 6 characters
  - **Status Codes**: 201 (Created), 400 (Bad Request), 409 (Conflict), 500 (Server Error)

- **POST** `/auth/login` - Authenticate user
  - **Body**: `{ "emailOrUsername": string, "password": string }`
  - **Returns**: User object + JWT token
  - **Accepts**: Either email or username as identifier
  - **Status Codes**: 200 (OK), 400 (Bad Request), 401 (Unauthorized), 500 (Server Error)

- **GET** `/auth/me` - Get current user profile
  - **Headers**: `Authorization: Bearer <token>`
  - **Returns**: User profile
  - **Status Codes**: 200 (OK), 401 (Unauthorized), 404 (Not Found), 500 (Server Error)

### Response Format

All endpoints return JSON in the following format:

**Success Response:**
```json
{
  "success": true,
  "data": { ... }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "error message"
}
```

### OpenAPI Documentation

View the complete API specification in `openapi.yaml` or use tools like Swagger UI:

```bash
# Using npx
npx @redocly/cli preview-docs openapi.yaml

# Or with Swagger Editor online
# https://editor.swagger.io/
```

## 🐳 Docker

### Build Image

```bash
docker-compose build
```

### Start Services

```bash
docker-compose up -d
```

### View Logs

```bash
docker-compose logs -f api
```

### Stop Services

```bash
docker-compose down
```

### Remove Volumes

```bash
docker-compose down -v
```

## 🔧 Configuration

Environment variables can be configured in `.env`:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=microservice_db

# Server
SERVER_PORT=8080
```

## 🔧 Integration with Other Services

This authentication service is designed to work as part of a microservices architecture:

### Token-Based Authentication Flow

1. **User Registration/Login**: Client calls `/auth/register` or `/auth/login`
2. **Token Issued**: Service returns JWT token valid for 24 hours
3. **Authenticated Requests**: Client includes token in `Authorization: Bearer <token>` header
4. **Service Validation**: Other microservices can validate tokens by:
   - Sharing the JWT secret
   - Calling `/auth/me` endpoint to verify token and get user info
   - Implementing their own JWT validation using the same secret

### Example: Validating Tokens in Other Services

```go
// Other microservices can validate tokens using the same JWT secret
import "github.com/golang-jwt/jwt/v5"

func validateToken(tokenString string) (int, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte("your-secret-key-change-this-in-production"), nil
    })
    // ... validation logic
}
```

### Environment Variables

Make sure to configure the same JWT secret across all services that need to validate tokens.

## 📊 Database Schema

The authentication service uses PostgreSQL with the following schema:

### Users Table

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    login VARCHAR(255) UNIQUE NOT NULL,      -- Username
    hashed_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
```

### Schema Management

Database schema is managed through SQL files in the `db/init/` directory. PostgreSQL automatically executes these files in alphabetical order when the container is first created.

**Note**: These SQL files only run on first database initialization. To reset the database, run:
```bash
docker-compose down -v  # Remove volumes
docker-compose up -d    # Recreate with fresh database
```

## 🧪 Testing

```bash
# Run tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔐 Security Features

This authentication service implements industry-standard security practices:

### Password Security
- **bcrypt hashing**: Passwords are hashed using bcrypt with default cost (10)
- **Never exposed**: Password hashes are never returned in API responses
- **Minimum length**: 6 characters required

### JWT Token Security
- **HS256 signing**: Tokens signed with HMAC-SHA256
- **24-hour expiry**: Tokens automatically expire after 24 hours
- **Secure secret**: JWT secret key (⚠️ change `your-secret-key-change-this-in-production` in production)

### Input Validation
- **Email format**: Validated using regex pattern
- **Username constraints**: 3-20 alphanumeric characters (including underscore)
- **Uniqueness checks**: Email and username must be unique
- **SQL injection protection**: Parameterized queries prevent SQL injection

### Best Practices for Production

- ⚠️ **Update JWT secret**: Change the hardcoded JWT secret in `app/internal/auth/jwt.go` to a strong random value
- 🔒 **Never commit `.env` files**: Keep sensitive data out of version control
- 🌐 **Enable HTTPS**: Always use TLS in production
- 🔑 **Use environment variables**: Store JWT secret and other sensitive data as environment variables
- 📦 **Keep dependencies updated**: Regularly update Go modules for security patches
- 🚫 **Rate limiting**: Consider adding rate limiting to prevent brute force attacks
- 📝 **Audit logging**: Log authentication attempts for security monitoring

## 📈 Performance

This template is optimized for:
- **Low memory footprint**: ~10-20MB per service
- **Fast startup**: <10ms
- **High concurrency**: Native Go goroutines
- **Efficient DB connections**: Connection pooling

## 🔄 CI/CD

The template is ready for CI/CD integration:

```yaml
# Example GitHub Actions workflow
- name: Build
  run: make build

- name: Test
  run: make test

- name: Docker Build
  run: make docker-build
```

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [OpenAPI Specification](https://swagger.io/specification/)

## 🤝 Contributing

This template is part of the AREA project. Follow the project's contribution guidelines.

## 📄 License

MIT License - See main project LICENSE file.

## 📝 Dependencies

- **github.com/lib/pq**: PostgreSQL driver
- **github.com/golang-jwt/jwt/v5**: JWT token generation and validation
- **golang.org/x/crypto**: bcrypt password hashing

---

**Authentication Service for AREA** - Automation platform inspired by IFTTT and Zapier
