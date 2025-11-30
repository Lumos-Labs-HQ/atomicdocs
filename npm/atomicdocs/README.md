<div align="center">
  <h1>⚛️ AtomicDocs</h1>
  <p><strong>Zero-config, auto-generated API documentation for Express.js and Hono</strong></p>
  <p>Built with Go and fasthttp for extreme performance</p>

  <p>
    <a href="https://www.npmjs.com/package/atomicdocs"><img src="https://img.shields.io/npm/v/atomicdocs.svg?style=flat-square" alt="npm version" /></a>
    <a href="https://www.npmjs.com/package/atomicdocs"><img src="https://img.shields.io/npm/dm/atomicdocs.svg?style=flat-square" alt="npm downloads" /></a>
    <a href="https://github.com/Lumos-Labs-HQ/atomicdocs/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-apache.svg?style=flat-square" alt="license" /></a>
  </p>
</div>

---

## ✨ Features

- 🚀 **Zero-config** — One line of middleware, no manual route definitions
- ⚡ **Blazing Fast** — Go server with fasthttp for documentation
- 🎯 **Interactive UI** — Swagger UI with "Try it out" functionality
- 🔄 **Auto-detection** — Automatically detects Express.js or Hono
- 📝 **Schema Extraction** — Parses request/response types from your code
- 📦 **Lightweight** — Minimal dependencies

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

> **Note:** No Go installation required! The binary is automatically downloaded during installation.

---

## 🚀 Quick Start

### Express.js

```javascript
const express = require('express');
const atomicdocs = require('atomicdocs');

const app = express();
app.use(express.json());

// Add AtomicDocs middleware
app.use(atomicdocs());

// Define your routes
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
  console.log(`🚀 Server: http://localhost:${PORT}`);
  console.log(`📚 Docs: http://localhost:${PORT}/docs`);
  
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

// Add AtomicDocs middleware (requires app instance and port)
app.use('*', atomicdocs(app, PORT));

// Define your routes
app.get('/users', (c) => c.json([{ id: 1, name: 'Alice' }]));

app.post('/users', async (c) => {
  const { name, email } = await c.req.json();
  return c.json({ id: 2, name, email }, 201);
});

app.get('/users/:id', (c) => {
  return c.json({ id: c.req.param('id'), name: 'Alice' });
});

serve({ fetch: app.fetch, port: PORT }, () => {
  console.log(`🚀 Server: http://localhost:${PORT}`);
  console.log(`📚 Docs: http://localhost:${PORT}/docs`);
});
```

### View Documentation

Once your server is running:

- **Swagger UI:** `http://localhost:<PORT>/docs`
- **OpenAPI JSON:** `http://localhost:<PORT>/docs/json`

---

## 📖 API Reference

### `atomicdocs()`

Express.js middleware. Auto-starts the documentation server.

```javascript
app.use(atomicdocs());
```

### `atomicdocs(app, port)`

Hono middleware. Requires app instance and port number.

```typescript
app.use('*', atomicdocs(app, 3000));
```

### `atomicdocs.register(app, port)`

Manually register Express routes. Call after all routes are defined.

```javascript
app.listen(3000, () => {
  atomicdocs.register(app, 3000);
});
```

---

## 🔍 What Gets Documented

AtomicDocs automatically detects and documents:

| Feature | Express | Hono |
|---------|---------|------|
| Route paths | ✅ | ✅ |
| HTTP methods (GET, POST, PUT, DELETE, PATCH) | ✅ | ✅ |
| Path parameters (`:id`, `:userId`) | ✅ | ✅ |
| Request body fields | ✅ | ✅ |
| Response codes | ✅ | ✅ |

### Schema Detection

```javascript
// Express - automatically extracts: name, email, age
app.post('/users', (req, res) => {
  const { name, email, age } = req.body;
  // ...
});

// Hono - automatically extracts: name, email, age
app.post('/users', async (c) => {
  const { name, email, age } = await c.req.json();
  // ...
});
```

### Smart Type Inference

| Field Name | Inferred Type | Example Value |
|------------|---------------|---------------|
| `age`, `count`, `id` | integer | `25` |
| `price`, `amount` | number | `99.99` |
| `email` | string | `user@example.com` |
| `isActive`, `enabled` | boolean | `true` |
| `name`, `title` | string | `"string"` |

---

## 🏗️ How It Works

```
┌─────────────────────────────┐
│     Your Express/Hono App   │
│                             │
│  app.use(atomicdocs())      │
└──────────────┬──────────────┘
               │ extracts routes
               ▼
┌─────────────────────────────┐
│   AtomicDocs Go Server      │
│   (runs on port 6174)       │
│                             │
│  • Receives route info      │
│  • Generates OpenAPI spec   │
│  • Serves Swagger UI        │
└──────────────┬──────────────┘
               │
               ▼
┌─────────────────────────────┐
│       Swagger UI            │
│    /docs endpoint           │
│                             │
│  • Interactive docs         │
│  • "Try it out" testing     │
└─────────────────────────────┘
```

---

## 💻 Supported Platforms

| Platform | Architecture |
|----------|-------------|
| Windows  | x64, arm64  |
| macOS    | x64, arm64 (M1/M2/M3) |
| Linux    | x64, arm64  |

---

## 🔧 TypeScript Support

Full TypeScript support with type definitions included:

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
  atomicdocs.register(app, PORT);
});
```

---

## 🐛 Troubleshooting

### Binary not found

```bash
# Re-run the install script
node node_modules/atomicdocs/install.js
```

### Permission denied (Linux/macOS)

```bash
chmod +x node_modules/atomicdocs/bin/atomicdocs-*
```

### Port 6174 already in use

```bash
# Find and kill the process
lsof -i :6174
kill -9 <PID>
```

### Docs page not loading

1. Make sure your server is running
2. Verify routes are registered: `atomicdocs.register(app, PORT)`
3. Check the console for errors

---

## 📁 Examples

Check out the [examples](https://github.com/Lumos-Labs-HQ/atomicdocs/tree/main/examples) directory:

- **express-demo** — Full Express.js example
- **hono-demo** — Full Hono example

---

## 🤝 Contributing

Contributions welcome! See the [GitHub repository](https://github.com/Lumos-Labs-HQ/atomicdocs).

---

## 📝 License

Apache 2.0 © [Lumos Labs HQ](https://github.com/Lumos-Labs-HQ)

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/Lumos-Labs-HQ">Lumos Labs HQ</a></p>
  <p>
    <a href="https://github.com/Lumos-Labs-HQ/atomicdocs">⭐ Star on GitHub</a>
    ·
    <a href="https://github.com/Lumos-Labs-HQ/atomicdocs/issues">Report Bug</a>
    ·
    <a href="https://github.com/Lumos-Labs-HQ/atomicdocs/issues">Request Feature</a>
  </p>
</div>
