# Microservice Template

A production-ready microservice template for the AREA project using Go, PostgreSQL, Docker, and OpenAPI.

## 🏗️ Architecture

This template provides a complete microservice setup with:

- **Go 1.22**: High-performance backend with native concurrency
- **PostgreSQL 16**: Robust relational database with ACID guarantees
- **Docker & Docker Compose**: Containerized deployment
- **OpenAPI 3.0**: API specification and documentation
- **RESTful API**: Clean HTTP endpoints with JSON responses

## 📁 Project Structure

```
.
├── app/                 # Application source code
│   ├── main.go         # Main application entry point
│   ├── go.mod          # Go module dependencies
│   └── go.sum          # Go dependency checksums
├── db/                  # Database configuration
│   └── init/           # Database initialization SQL files
│       ├── 01_create_tables.sql
│       └── 02_seed_data.sql
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile           # Multi-stage Docker build
├── Makefile             # Common commands
├── openapi.yaml         # OpenAPI 3.0 specification
├── .env.example         # Environment variables template
├── .gitignore           # Git ignore patterns
└── README.md            # This file
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

# Create a user
curl -X POST http://localhost:8080/users/create \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","first_name":"John","last_name":"Doe"}'

# Get all users
curl http://localhost:8080/users
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

### Users
- **GET** `/users` - Get all users
- **POST** `/users/create` - Create a new user

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

## 🏗️ Using This Template

### 1. Copy the Template

```bash
cp -r Backend/Template/Microservice Backend/MyNewService
cd Backend/MyNewService
```

### 2. Customize

- Update `app/go.mod` module name
- Modify `openapi.yaml` for your API spec
- Extend `app/main.go` with your business logic
- Add database schema files in `db/init/` directory
- Adjust environment variables in `.env.example`

### 3. Integrate

Each microservice runs independently with its own:
- Database instance
- API endpoints
- Docker container
- Network isolation

Services can communicate via Docker network or API gateway.

## 📊 Database Schema

Database schema is managed through SQL files in the `db/init/` directory. PostgreSQL automatically executes these files in alphabetical order when the container is first created.

### Adding Schema Files

1. **Create tables**: `db/init/01_create_tables.sql`
   ```sql
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       email VARCHAR(255) UNIQUE NOT NULL,
       first_name VARCHAR(255) NOT NULL,
       last_name VARCHAR(255) NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```

2. **Add indexes**: In the same file or separate files
   ```sql
   CREATE INDEX idx_users_email ON users(email);
   ```

3. **Seed data**: `db/init/02_seed_data.sql`
   ```sql
   INSERT INTO users (email, first_name, last_name)
   VALUES ('user@example.com', 'John', 'Doe');
   ```

### Naming Convention

Use numbered prefixes to control execution order:
- `01_create_tables.sql` - Table definitions
- `02_seed_data.sql` - Initial data
- `03_create_views.sql` - Views and functions
- etc.

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

## 🔐 Security Best Practices

- Never commit `.env` files
- Use environment variables for secrets
- Enable HTTPS in production
- Implement authentication/authorization
- Validate all inputs
- Use prepared statements (already implemented)
- Keep dependencies updated

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

---

**Built for AREA** - Automation platform inspired by IFTTT and Zapier
