# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ShowMeCode (FanCode) is a visual OJ (Online Judge) system - an educational platform combining code execution, debugging visualization, and a knowledge base for programming learning. It's a monorepo with 4 projects:

- **backend/** - Go backend API server (Gin + GORM + MySQL + Redis)
- **front-admin/** - Vue 3 admin dashboard for system management
- **front-user/** - Vue 3 user-facing coding platform with Monaco editor
- **go-debugger/** - Go DAP (Debug Adapter Protocol) server wrapping GDB for C/C++ debugging

## Build & Run Commands

### Backend (Go)
```bash
cd backend
go mod download
go build -o fancode .
./fancode  # Requires Linux (uses cgroups for sandboxing)

# Docker
docker build -t showmecode -f Dockerfile .
docker run -d -v /var/run/docker.sock:/var/run/docker.sock -v /var/fanCode:/var/fanCode --network=host showmecode

# Tests
go test ./...
go test -v ./service/system_service/...  # Run specific package tests
```

### Front-Admin (Vue 3 + Vite)
```bash
cd front-admin
pnpm install
pnpm dev              # Development server
pnpm build:pro        # Production build (includes vue-tsc type check)
pnpm lint             # ESLint
pnpm format           # Prettier
pnpm lint:style       # StyleLint
```

### Front-User (Vue 3 + Vite)
```bash
cd front-user
pnpm install
pnpm dev              # Development server
pnpm build:pro        # Production build
pnpm code:check       # lint + format check + style check
pnpm code:fix         # lint fix + format + style fix
```

### Go-Debugger
```bash
cd go-debugger
go mod download
go build -o godebugger .
./godebugger -port 8889 -file <executable> -codeFile <source> -language c

# Docker
docker build -t go-debugger -f Dockerfile .
```

## Architecture

### System Communication Flow
```
Front-Admin ──┐
              │ HTTP/REST
              ▼
Front-User ──► Backend (Go) ──► MySQL/Redis
              │     │
              │     │ Docker + DAP/TCP
              │     ▼
              └──► Go-Debugger (GDB wrapper)
```

### Backend Structure (Layered Architecture)
- **controller/** - HTTP handlers (admin/, user/ subdirectories)
- **service/** - Business logic layer
- **dao/** - Data access objects (database queries)
- **models/** - Data models (dto/ for transfer, po/ for persistence, vo/ for response)
- **routers/** - Route definitions
- **interceptor/** - Middleware (CORS, logging, rate limiting, auth)
- **common/** - Shared utilities (config, AI providers, file storage, error handling)
- Uses Google Wire for dependency injection (wire.go, wire_gen.go)

### Frontend Structure (Both front-admin and front-user)
- **src/api/** - API client functions organized by domain
- **src/views/** - Page components
- **src/components/** - Reusable components
- **src/store/modules/** - Pinia state stores
- **src/router/** - Vue Router configuration
- **src/utils/** - Utility functions

### Key Frontend Components (front-user)
- **components/code-editor/** - Monaco editor integration with debugging support
- **components/code-visual/** - Data structure visualization (@antv/g6)
- Uses xterm for terminal output, WebSocket for real-time debugging

### Go-Debugger Structure
- **debugger/** - Core debugger implementations
  - **c_debugger/** - C language debugger
  - **cpp_debugger/** - C++ language debugger
  - **gdb_debugger/gdb/** - GDB command parsing and grammar
- **protocol/** - DAP protocol request/response/event handling
- **server.go** - DAP server over TCP

## Key Technical Details

- Backend requires Linux (uses cgroups for code execution sandboxing)
- Backend uses Docker-in-Docker for isolated code execution
- Go-Debugger wraps GDB for C/C++ debugging via DAP protocol
- Frontend uses pnpm as package manager (enforced via preinstall script)
- Both frontends use Vue 3 Composition API with TypeScript
- Database schema in backend/fan_code.sql

## API Patterns

Backend API routes:
- `/api/auth/*` - Authentication
- `/api/account/*` - User accounts
- `/api/visual-document/*` - Learning documents
- `/api/saved-code/*` - User saved code
- `/api/debug/*` - Debugging endpoints
- `/manage/*` - Admin endpoints (api, role, user, menu management)
