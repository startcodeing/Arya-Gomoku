# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Arya-Gomoku** is a sophisticated Gomoku (Five-in-a-Row) game platform built with Golang backend and Vue 3 frontend. The project supports multiple game modes including AI vs Human, LLM vs Human, and Player vs Player battles with real-time WebSocket communication and persistent data storage.

## Development Commands

### Backend (Go)
```bash
cd backend

# Development
make dev          # Run with hot-reload (requires Air)
make run          # Build and run application
make build        # Build application
make test         # Run tests with verbose output
make fmt          # Format code
make lint         # Run linter (requires golangci-lint)

# Database
make migrate-up    # Apply all pending migrations
make migrate-down  # Rollback last migration
make migrate-status # Show migration status
make migrate-reset  # Reset database (WARNING: drops all data)

# Production
make build-prod   # Cross-platform production build
make docker-build  # Build Docker image
make docker-run   # Run Docker container

# Testing
make benchmark    # Run AI performance benchmark
```

### Frontend (Vue 3)
```bash
cd frontend

# Development
npm run dev       # Development server with hot-reload
npm run build     # Production build
npm run preview   # Preview production build

# Testing
npm run test      # Run tests (when test suite is added)
```

## Architecture Overview

### Backend Architecture (Golang)
- **Layered Architecture**: Controller → Service → Repository → Database
- **Web Framework**: Gin for RESTful API and WebSocket support
- **Database**: PostgreSQL (primary) with SQLite fallback for development
- **ORM**: GORM for database operations with repository pattern
- **Caching**: Redis for session management and performance optimization
- **Authentication**: JWT with refresh token support
- **Real-time**: WebSocket hub for live game state synchronization
- **Configuration**: Environment-based configuration with godotenv

### Frontend Architecture (Vue 3)
- **Framework**: Vue 3 with Composition API and `<script setup>`
- **Build Tool**: Vite for fast development and optimized builds
- **State Management**: Pinia for reactive global state
- **Routing**: Vue Router with protected route navigation
- **HTTP Client**: Axios with TypeScript interface definitions
- **UI Components**: Lucide Vue Next for icons
- **Proxy Configuration**: Vite proxies API calls to backend (localhost:8080)

## Key Technologies & Dependencies

### Backend Dependencies
- **Gin v1.9.1**: Web framework and routing
- **GORM v1.31.0**: ORM and database operations
- **PostgreSQL driver**: Database connectivity
- **gorilla/websocket**: WebSocket implementation
- **golang-jwt/jwt v5**: JWT authentication
- **redis**: Redis client for caching
- **godotenv**: Environment variable management

### Frontend Dependencies
- **Vue 3.4.0**: Frontend framework
- **TypeScript 5.2.0**: Type safety and developer experience
- **Vite 5.0.0**: Build tool and development server
- **Pinia 2.1.0**: State management
- **Vue Router 4.2.0**: Client-side routing
- **Axios 1.6.0**: HTTP client
- **Lucide Vue Next**: Icon library

## Game Modes & Implementation

### AI vs Human Mode
- **Classic AI**: Heuristic-based algorithm with threat detection
- **Enhanced AI**: Minimax with Alpha-Beta Pruning and transposition table
- **Difficulty Levels**: Easy, Medium, Hard, Expert (search depth 1-6)
- **Performance Tracking**: Search nodes, pruning efficiency, response time

### LLM vs Human Mode
- **LLM Integration**: Support for multiple LLM providers (OpenAI, Anthropic, etc.)
- **Dynamic Configuration**: Environment-based API configuration
- **Reasoning Display**: AI move reasoning and confidence scores

### Player vs Player Mode
- **Room System**: WebSocket-based real-time gameplay
- **Invitation System**: Shareable room links
- **Matchmaking**: Active room discovery and joining

### User & Persistence System
- **Authentication**: JWT with refresh token rotation
- **User Profiles**: Avatar, statistics, and game history
- **Data Persistence**: PostgreSQL with complete game state storage
- **Admin Features**: System statistics and management

## Configuration Setup

### Backend Environment (.env)
```bash
# Server Settings
SERVER_PORT=8080
ENVIRONMENT=development

# Database Settings (PostgreSQL recommended)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_NAME=gomoku
DB_SSL_MODE=disable

# Redis Settings (optional but recommended)
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT Settings
JWT_SECRET_KEY=your-super-secret-jwt-key
JWT_ACCESS_TOKEN_TTL=24h
```

### Frontend Development
- Development server runs on `localhost:5173`
- API calls proxied to `localhost:8080/api`
- TypeScript strict mode enabled for type safety

## Database Schema

The project uses a migration-based schema with the following main entities:
- **Users**: Authentication and profile data
- **Games**: Game sessions and metadata
- **Moves**: Individual move history with timestamps
- **Statistics**: User performance metrics

Run migrations with `make migrate-up` to set up the database.

## Key Files & Directories

### Backend Structure
```
backend/
├── main.go                 # Application entry point
├── internal/
│   ├── controller/         # HTTP handlers
│   ├── service/           # Business logic
│   ├── repository/        # Data access layer
│   ├── model/            # Data models
│   ├── database/         # Database configuration
│   └── middleware/       # HTTP middleware
├── migrations/            # Database migrations
├── cmd/                  # CLI tools and demos
└── .air.toml            # Hot-reload configuration
```

### Frontend Structure
```
frontend/
├── src/
│   ├── components/      # Vue components
│   ├── services/        # API service classes
│   ├── stores/          # Pinia state stores
│   ├── router/          # Vue Router configuration
│   ├── types/           # TypeScript definitions
│   └── utils/           # Utility functions
└── vite.config.ts       # Vite configuration
```

## Development Guidelines

### Code Standards
- **Go**: Follow Go official formatting conventions
- **TypeScript**: Use strict typing with complete annotations
- **Vue**: Use Composition API with `<script setup>` syntax
- **Testing**: Maintain test coverage for critical components

### Security Considerations
- **Authentication**: Always validate JWT tokens
- **CORS**: Configured for specific frontend origins
- **Input Validation**: Validate all request bodies
- **Database**: Use parameterized queries to prevent injection

### Performance Optimizations
- **Backend**: Connection pooling, Redis caching, AI algorithm optimizations
- **Frontend**: Component lazy loading, optimized bundle size, efficient state management

## WebSocket Implementation

The game uses WebSocket for real-time communication:
- **Connection**: `ws://localhost:8080/ws`
- **Hub Architecture**: Centralized WebSocket hub for room management
- **Message Types**: Game state updates, move synchronization, chat messages

## AI Algorithm Implementation

The AI engine uses multiple strategies:
1. **Threat Detection**: Identify immediate winning/defensive moves
2. **Position Scoring**: Evaluate board positions based on patterns
3. **Minimax Search**: Lookahead search with Alpha-Beta Pruning
4. **Transposition Table**: Cache board positions to avoid redundant calculations

## Troubleshooting

### Common Issues
1. **CORS Errors**: Verify frontend origin in CORS middleware
2. **Database Connection**: Check PostgreSQL service and credentials
3. **WebSocket Issues**: Verify firewall settings and port availability
4. **AI Performance**: Adjust search depth in AI configuration

### Debug Commands
```bash
# Backend debugging
go run main.go
make test              # Run specific tests
make migrate-status    # Check database status

# Frontend debugging
npm run build          # Check for build errors
npm run dev           # Development server with hot-reload
```

## Testing Strategy

### Backend Testing
- **Unit Tests**: Business logic and AI algorithm testing
- **Integration Tests**: API endpoint testing
- **Benchmark Tests**: AI performance measurement

### Frontend Testing
- **Type Checking**: TypeScript validation
- **Build Testing**: Production build verification
- **Component Testing**: Vue component unit tests (setup ready)

## Deployment

### Backend Deployment
```bash
# Production build
make build-prod

# Docker deployment
make docker-build && make docker-run
```

### Frontend Deployment
```bash
npm run build
# Deploy dist/ directory to static web server
```

The project is designed to run in a containerized environment with proper environment variable configuration for different deployment scenarios.