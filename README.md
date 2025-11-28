# AtomicDocs

Lightweight, auto-generated API documentation for Express.js and Hono. Built with Go and fasthttp for extreme performance.

## Features

- 🚀 **Zero-config** - Just add one line of middleware
- ⚡ **Fast** - Built with Go and fasthttp
- 🎯 **Interactive** - Full Swagger UI with "Try it out"
- 🔌 **Pluggable** - Extensible parser system for any framework
- 📦 **Lightweight** - Minimal dependencies

## Quick Start

### 1. Start AtomicDocs Server

```bash
go run cmd/server/main.go
```

Server runs on `http://localhost:6174`

### 2. Add to Your App

**Express.js:**
```javascript
const atomicDocs = require('./parsers/express');
app.use(atomicDocs(app));
```

**Hono:**
```typescript
import { atomicDocs } from './parsers/hono';
app.use('*', atomicDocs(app));
```

### 3. View Docs

Visit `http://localhost:6174/docs`

## Architecture

```
User App (Express/Hono)
    ↓ POST routes to /api/register
Go Server (fasthttp)
    ↓ Generates OpenAPI 3.0 spec
Swagger UI (/docs)
    ↓ Try it out → Direct API calls
User App (handles requests)
```

## API Endpoints

- `GET /docs` - Swagger UI
- `GET /docs/json` - OpenAPI 3.0 spec
- `POST /api/register` - Register routes (internal)

## Plugin System

Future parsers (Rust, Python, C++) will implement:

```go
type RouteInfo struct {
    Method      string
    Path        string
    Summary     string
    Description string
    // ... OpenAPI fields
}
```

Send via HTTP POST to `/api/register` with JSON array of routes.

## Examples

See `examples/express-demo` and `examples/hono-demo`

## License

MIT
