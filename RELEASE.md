# AtomicDocs v1.0.0 - Final Release

> 🎉 **This is the stable v1.0.0 release of AtomicDocs!**

## Overview

AtomicDocs is a zero-config, auto-generated API documentation tool for Express.js and Hono frameworks. Built with Go and fasthttp for extreme performance.

## Installation

```bash
npm install atomicdocs
```

That's it! No Go installation required. The binary is automatically downloaded during installation.

## Quick Start

### Express.js

```javascript
const express = require('express');
const atomicdocs = require('atomicdocs');

const app = express();
app.use(express.json());
app.use(atomicdocs());

app.get('/users', (req, res) => res.json([{ id: 1, name: 'John' }]));
app.post('/users', (req, res) => res.status(201).json({ id: 2, ...req.body }));

const PORT = 3000;
app.listen(PORT, () => {
  console.log(`Server: http://localhost:${PORT}`);
  console.log(`Docs: http://localhost:${PORT}/docs`);
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

app.use('*', atomicdocs(app, PORT));

app.get('/users', (c) => c.json([{ id: 1, name: 'Alice' }]));
app.post('/users', async (c) => c.json({ id: 2, ...(await c.req.json()) }, 201));

serve({ fetch: app.fetch, port: PORT }, () => {
  console.log(`Server: http://localhost:${PORT}`);
  console.log(`Docs: http://localhost:${PORT}/docs`);
});
```

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│              npm install atomicdocs                         │
└──────────────────────────┬──────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│    postinstall downloads binary from GitHub Release         │
│    Binary saved to node_modules/atomicdocs/bin/             │
└──────────────────────────┬──────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│    JavaScript auto-detects OS/architecture                  │
│    Spawns correct binary (Win/macOS/Linux, x64/arm64)       │
└──────────────────────────┬──────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│    Go server runs on localhost:6174                         │
│    Express/Hono app sends routes to Go server               │
└──────────────────────────┬──────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│    Go analyzes code and generates OpenAPI 3.0 spec          │
│    Swagger UI served at /docs                               │
└─────────────────────────────────────────────────────────────┘
```

---

## Supported Platforms

| Platform | Architecture | Binary |
|----------|-------------|--------|
| Windows  | x64         | `atomicdocs-win-x64.exe` |
| Windows  | arm64       | `atomicdocs-win-arm64.exe` |
| macOS    | x64 (Intel) | `atomicdocs-darwin-x64` |
| macOS    | arm64 (M1/M2/M3) | `atomicdocs-darwin-arm64` |
| Linux    | x64         | `atomicdocs-linux-x64` |
| Linux    | arm64       | `atomicdocs-linux-arm64` |

---

## Features

### ✅ Core Features

| Feature | Description |
|---------|-------------|
| 🚀 Zero-config | Just `app.use(atomicdocs())` - no manual route definitions |
| 🔍 Auto route detection | Discovers all GET/POST/PUT/DELETE/PATCH routes automatically |
| 📝 Schema extraction | Analyzes handler code for request/response parameters |
| 📋 OpenAPI 3.0 | Generates valid OpenAPI 3.0 JSON specification |
| 🎯 Swagger UI | Interactive UI with "Try it out" functionality |
| 🧪 Real-time testing | Test APIs directly from browser |
| 💻 Cross-platform | Windows, macOS, Linux (x64 & arm64) |
| ⚡ Blazing fast | Built with Go and fasthttp |
| 📦 Lightweight | Minimal memory footprint |

### ✅ Schema Detection

| Detection | Example |
|-----------|---------|
| Express body | `const { x, y } = req.body` |
| Hono body | `const { x, y } = await c.req.json()` |
| Path parameters | `:id`, `:userId` |
| Type inference | `age` → integer, `price` → number |
| Example values | `email` → `user@example.com` |

### ✅ HTTP Methods Supported

- `GET` - List and retrieve resources
- `POST` - Create new resources
- `PUT` - Update existing resources
- `DELETE` - Remove resources
- `PATCH` - Partial updates

### ✅ Response Codes Documented

- `200` - OK
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found

---

## Package Structure

```
atomicdocs/
├── index.js              # Main entry - auto-detects Express/Hono
├── index.d.ts            # TypeScript definitions
├── install.js            # Downloads binary from GitHub Release
├── package.json
└── bin/
    └── atomicdocs-*      # Platform-specific binary (downloaded)
```

---

## Troubleshooting

### Binary not found error

```
Error: Cannot find module '.../bin/atomicdocs-...'
```

**Solution:**
1. Check GitHub Release exists: `https://github.com/Lumos-Labs-HQ/atomicdocs/releases/tag/v1.0.0`
2. Re-run: `npm rebuild atomicdocs`
3. Or manually run: `node node_modules/atomicdocs/install.js`

### Permission denied (Linux/macOS)

```
Error: EACCES: permission denied
```

**Solution:**
```bash
chmod +x node_modules/atomicdocs/bin/atomicdocs-*
```

### Port 6174 already in use

```
Error: listen EADDRINUSE: address already in use :::6174
```

**Solution:**
```bash
# Find and kill existing process
lsof -i :6174
kill -9 <PID>
```

### Docs not loading

1. Ensure server is running
2. Check routes are registered: `atomicdocs.register(app, PORT)`
3. Visit `http://localhost:<PORT>/docs`

---

## CI/CD Workflows

### Release Workflow (`release.yml`)

**Trigger:** Git tag push (e.g., `git tag v1.0.0 && git push origin v1.0.0`)

**Steps:**
1. Checkout code
2. Setup Go 1.22
3. Build 6 binaries with `-ldflags="-s -w"` (stripped)
4. Create GitHub Release
5. Upload binaries as release assets

### NPM Publish Workflow (`npm-publish.yml`)

**Trigger:** After release workflow completes

**Steps:**
1. Get version from git tag
2. Update package.json version
3. Publish to NPM registry

---

## How to Release (For Maintainers)

```bash
# 1. Update version
cd npm/atomicdocs
npm version 1.0.0 --no-git-tag-version

# 2. Commit
git add .
git commit -m "Release v1.0.0"
git push origin main

# 3. Tag and push
git tag v1.0.0
git push origin v1.0.0

# 4. Verify (after CI completes)
npm view atomicdocs version
```

---

## GitHub Secrets Required

| Secret | Description | How to get |
|--------|-------------|-----------|
| `NPM_TOKEN` | NPM automation token | https://www.npmjs.com/settings/YOUR_USERNAME/tokens |

---

## v1.0.0 Release Notes

**Release Date:** November 2025

### 🎉 What's Included

- ✅ **Express.js support** - Full middleware integration
- ✅ **Hono support** - Full middleware integration
- ✅ **OpenAPI 3.0 generation** - Valid spec output
- ✅ **Swagger UI** - Interactive "Try it out" functionality
- ✅ **Auto schema extraction** - Parses handler code for parameters
- ✅ **Cross-platform binaries** - Windows, macOS, Linux (x64 & arm64)
- ✅ **Smart type inference** - Detects integer, number, boolean, string
- ✅ **Path parameter detection** - `:id`, `:userId`, etc.
- ✅ **Zod schema support** - Extracts types from Zod schemas
- ✅ **TypeScript support** - Full type definitions included

### 📊 Binary Sizes

| Platform | Size (approx) |
|----------|---------------|
| Windows x64 | ~8 MB |
| Windows arm64 | ~8 MB |
| macOS x64 | ~8 MB |
| macOS arm64 | ~8 MB |
| Linux x64 | ~8 MB |
| Linux arm64 | ~8 MB |

### 🔗 Links

- **NPM:** https://www.npmjs.com/package/atomicdocs
- **GitHub:** https://github.com/Lumos-Labs-HQ/atomicdocs
- **Issues:** https://github.com/Lumos-Labs-HQ/atomicdocs/issues

---

## License

Apache 2.0 © [Lumos Labs HQ](https://github.com/Lumos-Labs-HQ)

---

<div align="center">
  <h3>🎉 Thank you for using AtomicDocs!</h3>
  <p>If you find this project useful, please ⭐ star it on GitHub!</p>
</div>
