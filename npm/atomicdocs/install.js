const https = require('https');
const fs = require('fs');
const path = require('path');

const REPO = 'Lumos-Labs-HQ/atomicdocs';
const VERSION = require('./package.json').version;

function getBinaryName() {
  const platform = process.platform;
  const arch = process.arch;
  
  console.log(`Detecting platform: ${platform}, arch: ${arch}`);
  
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
    throw new Error(`Unsupported platform: ${platform}-${arch}. Supported: win32/darwin/linux with x64/arm64`);
  }
  
  const ext = platform === 'win32' ? '.exe' : '';
  return `atomicdocs-${mappedPlatform}-${mappedArch}${ext}`;
}

function followRedirects(url, maxRedirects = 5) {
  return new Promise((resolve, reject) => {
    if (maxRedirects === 0) {
      return reject(new Error('Too many redirects'));
    }
    
    https.get(url, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        const redirectUrl = response.headers.location;
        console.log(`  Redirecting to: ${redirectUrl.substring(0, 80)}...`);
        followRedirects(redirectUrl, maxRedirects - 1).then(resolve).catch(reject);
      } else if (response.statusCode === 200) {
        resolve(response);
      } else if (response.statusCode === 404) {
        reject(new Error(`Binary not found at ${url}. Make sure GitHub Release v${VERSION} exists with the binary uploaded.`));
      } else {
        reject(new Error(`HTTP ${response.statusCode}: ${response.statusMessage}`));
      }
    }).on('error', reject);
  });
}

function downloadBinary() {
  const binaryName = getBinaryName();
  const binDir = path.join(__dirname, 'bin');
  const binaryPath = path.join(binDir, binaryName);
  
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }
  
  if (fs.existsSync(binaryPath)) {
    const stats = fs.statSync(binaryPath);
    if (stats.size > 0) {
      console.log('✓ Binary already exists');
      return Promise.resolve();
    }
    // Remove empty/corrupted binary
    fs.unlinkSync(binaryPath);
  }
  
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${binaryName}`;
  
  console.log(`Downloading ${binaryName} for version ${VERSION}...`);
  console.log(`  URL: ${url}`);
  
  return followRedirects(url)
    .then((response) => {
      return new Promise((resolve, reject) => {
        const file = fs.createWriteStream(binaryPath);
        response.pipe(file);
        file.on('finish', () => {
          file.close();
          try {
            fs.chmodSync(binaryPath, '755');
          } catch (e) {
            // chmod may fail on Windows, that's ok
          }
          console.log('✓ Binary downloaded successfully');
          resolve();
        });
        file.on('error', (err) => {
          if (fs.existsSync(binaryPath)) fs.unlinkSync(binaryPath);
          reject(err);
        });
      });
    })
    .catch((err) => {
      if (fs.existsSync(binaryPath)) fs.unlinkSync(binaryPath);
      throw err;
    });
}

downloadBinary()
  .then(() => {
    console.log('✓ AtomicDocs installation complete');
  })
  .catch((err) => {
    console.error('\n❌ Installation failed:', err.message);
    console.error('\nTroubleshooting:');
    console.error('  1. Check if GitHub Release v' + VERSION + ' exists at:');
    console.error('     https://github.com/' + REPO + '/releases/tag/v' + VERSION);
    console.error('  2. Verify the binary assets are uploaded to the release');
    console.error('  3. Your platform: ' + process.platform + '-' + process.arch);
    process.exit(1);
  });
