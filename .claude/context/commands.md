# Commands Reference

## Quick Start Commands

### Development
```bash
# Install all dependencies (backend + frontend)
make install

# Start both services in development mode
make dev
# Backend: http://localhost:8080
# Frontend: http://localhost:5173

# Backend only
cd backend && go run cmd/main.go

# Frontend only
cd frontend && npm run dev
```

### Docker (Recommended for Testing)
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild after changes
docker-compose up -d --build
```

### Testing
```bash
# Run all tests
make test

# Backend tests only
cd backend && go test ./...

# Backend tests with coverage
cd backend && go test -cover ./...

# Frontend tests only
cd frontend && npm test
```

### Building
```bash
# Build everything
make build

# Backend build
cd backend && go build -o bin/ai-motion cmd/main.go

# Frontend build
cd frontend && npm run build
```
