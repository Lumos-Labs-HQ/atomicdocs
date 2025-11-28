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

function extractRoutes(app) {
  const routes = [];
  
  app._router.stack.forEach(layer => {
    if (layer.route) {
      const handler = layer.route.stack[0]?.handle;
      
      Object.keys(layer.route.methods).forEach(method => {
        if (layer.route.methods[method]) {
          routes.push({
            method: method.toUpperCase(),
            path: layer.route.path,
            handler: handler ? handler.toString() : ''
          });
        }
      });
    }
  });
  
  return routes.filter(r => !r.path.startsWith('/docs'));
}

function registerRoutes(routes, port) {
  const data = JSON.stringify({ routes, port });
  const req = http.request({
    hostname: 'localhost',
    port: 6174,
    path: '/api/register',
    method: 'POST',
    headers: { 
      'Content-Type': 'application/json',
      'Content-Length': data.length
    }
  });
  
  req.on('error', () => {});
  req.write(data);
  req.end();
}

module.exports = function() {
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
};

module.exports.register = function(app, port) {
  if (!serverReady) {
    setTimeout(() => module.exports.register(app, port), 100);
    return;
  }
  const routes = extractRoutes(app);
  registerRoutes(routes, port);
  console.log(`✓ Registered ${routes.length} routes with AtomicDocs`);
};
