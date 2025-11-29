const { spawn } = require('child_process');
const http = require('http');
const path = require('path');

let goServer = null;
let serverReady = false;

function getBinaryName() {
  const platform = process.platform;
  const arch = process.arch;
  
  const platformMap = {
    'win32': 'win',
    'darwin': 'darwin',
    'linux': 'linux'
  };
  
  const archMap = {
    'x64': 'x64',
    'arm64': 'arm64'
  };
  
  const mappedPlatform = platformMap[platform];
  const mappedArch = archMap[arch];
  
  if (!mappedPlatform || !mappedArch) {
    throw new Error(`Unsupported platform: ${platform}-${arch}`);
  }
  
  const ext = platform === 'win32' ? '.exe' : '';
  return `atomicdocs-${mappedPlatform}-${mappedArch}${ext}`;
}

function startGoServer() {
  if (goServer) return;
  
  const binaryName = getBinaryName();
  const binaryPath = path.join(__dirname, 'bin', binaryName);
  
  goServer = spawn(binaryPath, [], { 
    detached: false,
    stdio: ['ignore', 'pipe', 'pipe']
  });
  
  goServer.on('error', (err) => {
    console.error('Failed to start AtomicDocs:', err.message);
  });
  
  setTimeout(() => { serverReady = true; }, 500);
}

function isHono(app) {
  return app && app.routes && Array.isArray(app.routes);
}

function extractExpressRoutes(app) {
  const routes = [];
  const fs = require('fs');
  const path = require('path');
  
  function getFilePath(handler) {
    if (!handler) return null;
    
    try {
      // Get the source location from the function
      const funcStr = handler.toString();
      const lines = funcStr.split('\n');
      
      // Try to find file path from function's source
      const oldPrepare = Error.prepareStackTrace;
      Error.prepareStackTrace = (_, stack) => stack;
      const stack = new Error().stack;
      Error.prepareStackTrace = oldPrepare;
      
      // Look through call sites for .ts or .js files (not node_modules)
      for (let i = 0; i < stack.length; i++) {
        const fileName = stack[i].getFileName();
        if (fileName && 
            (fileName.endsWith('.ts') || fileName.endsWith('.js')) &&
            !fileName.includes('node_modules') &&
            !fileName.includes('atomicdocs')) {
          return fileName;
        }
      }
      
      // Fallback: try to parse from stack string
      const err = new Error();
      const stackLines = err.stack.split('\n');
      for (let line of stackLines) {
        const match = line.match(/\((.+\.(ts|js)):(\d+):(\d+)\)/);
        if (match && !match[1].includes('node_modules') && !match[1].includes('atomicdocs')) {
          return match[1];
        }
      }
    } catch (e) {
      console.error('[AtomicDocs] Error getting file path:', e.message);
    }
    return null;
  }
  
  function extractImports(filePath) {
    if (!filePath || !fs.existsSync(filePath)) return [];
    
    try {
      const content = fs.readFileSync(filePath, 'utf-8');
      const imports = [];
      
      // Match: import { X, Y } from './file'
      const importRegex = /import\s*\{([^}]+)\}\s*from\s*['"]([^'"]+)['"]/g;
      let match;
      while ((match = importRegex.exec(content)) !== null) {
        const names = match[1].split(',').map(n => n.trim());
        const from = match[2];
        names.forEach(name => {
          imports.push({ name, from });
        });
      }
      
      return imports;
    } catch (e) {
      console.error(`[AtomicDocs] Error reading file ${filePath}:`, e.message);
      return [];
    }
  }
  
  function extractFromStack(stack, basePath = '') {
    stack.forEach(layer => {
      if (layer.route) {
        const handler = layer.route.stack[0]?.handle;
        const handlerCode = handler ? handler.toString() : '';
        
        // Detect controller file from route path
        let filePath = null;
        const routePath = basePath + layer.route.path;
        const cwd = process.cwd();
        
        // Try to find controller file based on route
        const possiblePaths = [
          path.join(cwd, 'src/controllers/auth.controller.ts'),
          path.join(cwd, 'src/controllers/auth.controller.js'),
          path.join(cwd, 'src/controllers/user.controller.ts'),
          path.join(cwd, 'src/controllers/user.controller.js'),
          path.join(cwd, 'src/controllers/post.controller.ts'),
          path.join(cwd, 'src/controllers/post.controller.js'),
        ];
        
        // Match route to controller
        if (routePath.includes('/auth')) {
          filePath = possiblePaths.find(p => p.includes('auth') && fs.existsSync(p));
        } else if (routePath.includes('/user')) {
          filePath = possiblePaths.find(p => p.includes('user') && fs.existsSync(p));
        } else if (routePath.includes('/post')) {
          filePath = possiblePaths.find(p => p.includes('post') && fs.existsSync(p));
        }
        
        const imports = filePath ? extractImports(filePath) : [];
        
        Object.keys(layer.route.methods).forEach(method => {
          if (layer.route.methods[method]) {
            routes.push({
              method: method.toUpperCase(),
              path: basePath + layer.route.path,
              handler: handlerCode,
              filePath: filePath || 'unknown',
              imports: imports
            });
          }
        });
      } else if (layer.name === 'router' && layer.handle.stack) {
        const routePath = layer.regexp.source
          .replace('\\/?', '')
          .replace('(?=\\/|$)', '')
          .replace(/\\\//g, '/')
          .replace(/\^/g, '')
          .replace(/\$/g, '');
        extractFromStack(layer.handle.stack, routePath);
      }
    });
  }
  
  extractFromStack(app._router.stack);
  return routes.filter(r => !r.path.startsWith('/docs'));
}

function extractHonoRoutes(app) {
  return app.routes.map(route => ({
    method: route.method.toUpperCase(),
    path: route.path,
    handler: route.handler ? route.handler.toString() : ''
  })).filter(r => !r.path.startsWith('/docs'));
}

function registerRoutes(routes, port) {
  const fs = require('fs');
  const path = require('path');
  const schemaFiles = {};
  
  console.log('[AtomicDocs] Collecting schema files...');
  
  // Collect all schema files
  routes.forEach(route => {
    if (route.imports && route.filePath) {
      route.imports.forEach(imp => {
        const dir = path.dirname(route.filePath);
        let resolved = path.resolve(dir, imp.from);
        
        // Try extensions
        const extensions = ['.ts', '.js', '.tsx', '.jsx'];
        for (let ext of extensions) {
          if (fs.existsSync(resolved + ext)) {
            resolved = resolved + ext;
            break;
          }
        }
        
        if (fs.existsSync(resolved) && !schemaFiles[resolved]) {
          try {
            const content = fs.readFileSync(resolved, 'utf-8');
            schemaFiles[resolved] = content;
            console.log(`  ✓ Read schema file: ${resolved} (${content.length} bytes)`);
          } catch (e) {
            console.error(`  ✗ Failed to read schema file: ${resolved}`);
          }
        }
      });
    }
  });
  
  console.log(`[AtomicDocs] Sending ${Object.keys(schemaFiles).length} schema files to Go server`);
  
  const data = JSON.stringify({ 
    routes, 
    port,
    schemaFiles 
  });
  
  const req = http.request({
    hostname: 'localhost',
    port: 6174,
    path: '/api/register',
    method: 'POST',
    headers: { 
      'Content-Type': 'application/json',
      'Content-Length': Buffer.byteLength(data)
    }
  });
  
  req.on('error', (err) => {
    console.error('Failed to register routes:', err.message);
  });
  
  req.write(data);
  req.end();
}

// Express middleware
function expressMiddleware() {
  startGoServer();
  
  return function(req, res, next) {
    if (req.path === '/docs' || req.path === '/docs/json') {
      const options = {
        hostname: 'localhost',
        port: 6174,
        path: req.path,
        method: 'GET',
        headers: { 'X-App-Port': req.app.get('port') || req.socket.localPort }
      };
      
      http.get(options, (goRes) => {
        res.writeHead(goRes.statusCode, goRes.headers);
        goRes.pipe(res);
      }).on('error', (err) => {
        res.status(503).send('AtomicDocs unavailable');
      });
      return;
    }
    
    next();
  };
}

// Hono middleware
function honoMiddleware(app, port) {
  startGoServer();
  
  setTimeout(() => {
    if (!serverReady) return;
    const routes = extractHonoRoutes(app);
    registerRoutes(routes, port);
    console.log(`✓ Registered ${routes.length} routes with AtomicDocs`);
  }, 1000);
  
  return async (c, next) => {
    if (c.req.path === '/docs' || c.req.path === '/docs/json') {
      return new Promise((resolve) => {
        http.get({
          hostname: 'localhost',
          port: 6174,
          path: c.req.path,
          headers: { 'X-App-Port': port.toString() }
        }, (res) => {
          let body = '';
          res.on('data', chunk => body += chunk);
          res.on('end', () => {
            const contentType = res.headers['content-type'] || 'text/html';
            resolve(new Response(body, {
              status: res.statusCode,
              headers: { 'Content-Type': contentType }
            }));
          });
        }).on('error', () => {
          resolve(c.text('AtomicDocs unavailable', 503));
        });
      });
    }
    
    await next();
  };
}

// Auto-detect framework
module.exports = function(app, port) {
  if (isHono(app)) {
    return honoMiddleware(app, port);
  }
  return expressMiddleware();
};

// Express manual registration
module.exports.register = function(app, port) {
  if (!serverReady) {
    setTimeout(() => module.exports.register(app, port), 100);
    return;
  }
  const routes = extractExpressRoutes(app);
  console.log(`[AtomicDocs] Extracted ${routes.length} routes`);
  routes.forEach(r => {
    console.log(`  ${r.method} ${r.path} - File: ${r.filePath || 'unknown'}, Imports: ${r.imports?.length || 0}`);
  });
  registerRoutes(routes, port);
  console.log(`✓ Registered ${routes.length} routes with AtomicDocs`);
};
