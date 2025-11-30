<div align="center">
  <h1>⚛️ AtomicDocs</h1>
  <p><strong>Zero-config, auto-generated API documentation for Express.js and Hono</strong></p>
  <p>Built with Go and fasthttp for extreme performance</p>

  <p>
    <a href="https://www.npmjs.com/package/atomicdocs"><img src="https://img.shields.io/npm/v/atomicdocs.svg?style=flat-square" alt="npm version" /></a>
    <a href="https://www.npmjs.com/package/atomicdocs"><img src="https://img.shields.io/npm/dm/atomicdocs.svg?style=flat-square" alt="npm downloads" /></a>
    <a href="https://github.com/Lumos-Labs-HQ/atomicdocs/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg?style=flat-square" alt="license" /></a>
  </p>
</div>

---

## ✨ Features

- 🚀 **Zero-config** — Just add one line of middleware, no manual route definitions
- ⚡ **Lightning Fast** — Go server with fasthttp handles documentation at blazing speed
- 🎯 **Interactive UI** — Full Swagger UI with "Try it out" functionality
- 🔄 **Auto-detection** — Automatically detects Express.js or Hono framework
- 📝 **Schema Extraction** — Parses Zod/TypeScript schemas for request/response types
- 🔌 **Pluggable** — Extensible parser system for any framework
- 📦 **Lightweight** — Minimal dependencies, small footprint

---

## 📦 Installation

```bash
npm install atomicdocs
```

```bash
yarn add atomicdocs
```

```bash
pnpm add atomicdocs
```

The package automatically downloads the appropriate binary for your platform during installation.

### Supported Platforms

| Platform | Architecture |
|----------|-------------|
| Windows  | x64, arm64  |
| macOS    | x64, arm64  |
| Linux    | x64, arm64  |

---

## 🚀 Quick Start

### Express.js

```javascript
const express = require('express');
const atomicdocs = require('atomicdocs');

const app = express();
app.use(express.json());

// Add AtomicDocs middleware — that's it!
app.use(atomicdocs());

// Define your routes as usual
app.get('/users', (req, res) => {
  res.json([{ id: 1, name: 'John' }]);
});

app.post('/users', (req, res) => {
  const { name, email } = req.body;
  res.status(201).json({ id: 2, name, email });
});

app.get('/users/:id', (req, res) => {
  res.json({ id: req.params.id, name: 'John' });
});

const PORT = 3000;
app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
  console.log(`📚 Docs available at http://localhost:${PORT}/docs`);
  
  // Register routes after server starts
  atomicdocs.register(app, PORT);
});
```

### Hono

```typescript
import { Hono } from 'hono';
import { serve } from '@hono/node-server';
import atomicdocs from 'atomicdocs';

const app = new Hono();
const PORT = 3000;

// Add AtomicDocs middleware with app instance and port
app.use('*', atomicdocs(app, PORT));

// Define your routes as usual
app.get('/users', (c) => c.json([{ id: 1, name: 'Alice' }]));

app.post('/users', async (c) => {
  const { name, email } = await c.req.json();
  return c.json({ id: 2, name, email }, 201);
});

app.get('/users/:id', (c) => {
  const id = c.req.param('id');
  return c.json({ id, name: 'Alice' });
});

serve({ fetch: app.fetch, port: PORT }, () => {
  console.log(`Server running on http://localhost:${PORT}`);
  console.log(`📚 Docs available at http://localhost:${PORT}/docs`);
});
```

### View Documentation

Once your server is running, visit:

- **Swagger UI**: `http://localhost:<PORT>/docs`
- **OpenAPI JSON**: `http://localhost:<PORT>/docs/json`

---

## 📖 API Reference

### `atomicdocs()`

Returns Express.js middleware. Auto-starts the Go documentation server.

```javascript
app.use(atomicdocs());
```

### `atomicdocs(app, port)`

Returns Hono middleware. Requires app instance and port number.

```typescript
app.use('*', atomicdocs(app, 3000));
```

### `atomicdocs.register(app, port)`

Manually register Express routes. Call this after all routes are defined and the server has started.

```javascript
app.listen(3000, () => {
  atomicdocs.register(app, 3000);
});
```

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Your Application                        │
│                   (Express.js / Hono)                       │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           │ AtomicDocs middleware
                           │ extracts routes & schemas
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                   AtomicDocs Go Server                      │
│                  (fasthttp on port 6174)                    │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Route     │  │   OpenAPI   │  │     Swagger UI      │  │
│  │  Registry   │──│  Generator  │──│      Server         │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                      Swagger UI                             │
│                   /docs endpoint                            │
│                                                             │
│           "Try it out" → Direct API calls →                 │
│                 Your Application                            │
└─────────────────────────────────────────────────────────────┘
```

### How It Works

1. **Middleware Integration**: AtomicDocs middleware is added to your Express.js or Hono application
2. **Route Extraction**: The middleware automatically extracts all registered routes, including HTTP methods, paths, and handler information
3. **Schema Parsing**: If you're using Zod or TypeScript schemas, AtomicDocs parses them for request/response documentation
4. **Go Server**: A high-performance Go server (using fasthttp) runs on port `6174` to handle documentation generation
5. **OpenAPI Generation**: Routes are converted to OpenAPI 3.0 specification
6. **Swagger UI**: Interactive documentation is served at `/docs`

---

## 🔧 Advanced Usage

### TypeScript with Express

```typescript
import express, { Request, Response } from 'express';
import atomicdocs from 'atomicdocs';

const app = express();
app.use(express.json());
app.use(atomicdocs());

interface User {
  id: number;
  name: string;
  email: string;
}

app.get('/users', (req: Request, res: Response) => {
  const users: User[] = [{ id: 1, name: 'John', email: 'john@example.com' }];
  res.json(users);
});

const PORT = 3000;
app.listen(PORT, () => {
  console.log(`Docs at http://localhost:${PORT}/docs`);
  atomicdocs.register(app, PORT);
});
```

### With Zod Schema Validation

AtomicDocs automatically extracts Zod schemas for better documentation:

```typescript
// schemas/user.schema.ts
import { z } from 'zod';

export const CreateUserSchema = z.object({
  name: z.string().min(1),
  email: z.string().email(),
  age: z.number().optional(),
});

export const UserResponseSchema = z.object({
  id: z.number(),
  name: z.string(),
  email: z.string(),
  createdAt: z.string(),
});

// controllers/user.controller.ts
import { CreateUserSchema } from '../schemas/user.schema';

export const createUser = (req, res) => {
  const data = CreateUserSchema.parse(req.body);
  // ... create user logic
  res.status(201).json({ id: 1, ...data });
};
```

### Project Structure (Recommended)

```
src/
├── controllers/
│   ├── auth.controller.ts
│   ├── user.controller.ts
│   └── post.controller.ts
├── routes/
│   ├── auth.routes.ts
│   ├── user.routes.ts
│   └── post.routes.ts
├── schemas/
│   ├── auth.schemas.ts
│   ├── user.schemas.ts
│   └── post.schemas.ts
└── server.ts
```

AtomicDocs will automatically detect and document schemas from your `schemas/` directory.

---

## 📁 Examples

### Express.js Demo

```bash
cd examples/express-demo
npm install
node server.js
```

Then visit `http://localhost:3000/docs`

### Hono Demo

```bash
cd examples/hono-demo
npm install  # or bun install
npx tsx server.ts  # or bun run server.ts
```

Then visit `http://localhost:8080/docs`

---

## 🛠️ Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/Lumos-Labs-HQ/atomicdocs.git
cd atomicdocs

# Build the Go server
go build -o bin/atomicdocs cmd/server/main.go

# Run the server directly
./bin/atomicdocs
```

### Running the Go Server Standalone

```bash
go run cmd/server/main.go
```

The server runs on `http://localhost:6174` with these endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/docs` | GET | Swagger UI interface |
| `/docs/json` | GET | OpenAPI 3.0 JSON spec |
| `/api/register` | POST | Register routes (internal) |

---

## 🔌 Plugin System

AtomicDocs is designed to be extensible. Future parsers for other languages/frameworks will implement:

```go
type RouteInfo struct {
    Method      string            `json:"method"`
    Path        string            `json:"path"`
    Summary     string            `json:"summary"`
    Description string            `json:"description"`
    Tags        []string          `json:"tags"`
    Parameters  []Parameter       `json:"parameters"`
    RequestBody *RequestBody      `json:"requestBody"`
    Responses   map[string]Response `json:"responses"`
}
```

Send routes via HTTP POST to `/api/register` with JSON array.

### Planned Framework Support

- [ ] Fastify (Node.js)
- [ ] Koa (Node.js)
- [ ] Flask (Python)
- [ ] FastAPI (Python)
- [ ] Axum (Rust)
- [ ] Actix (Rust)

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

Apache 2.0 © [Lumos Labs HQ](https://github.com/Lumos-Labs-HQ)

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/Lumos-Labs-HQ">Lumos Labs HQ</a></p>
  <p>
    <a href="https://github.com/Lumos-Labs-HQ/atomicdocs/issues">Report Bug</a>
    ·
    <a href="https://github.com/Lumos-Labs-HQ/atomicdocs/issues">Request Feature</a>
  </p>
</div>
