const https = require('https');
const fs = require('fs');
const path = require('path');

const REPO = 'yourusername/atomicdocs';
const VERSION = require('./package.json').version;

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

function downloadBinary() {
  const binaryName = getBinaryName();
  const binDir = path.join(__dirname, 'bin');
  const binaryPath = path.join(binDir, binaryName);
  
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }
  
  if (fs.existsSync(binaryPath)) {
    console.log('✓ Binary already exists');
    return;
  }
  
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${binaryName}`;
  
  console.log(`Downloading ${binaryName}...`);
  
  const file = fs.createWriteStream(binaryPath);
  
  https.get(url, (response) => {
    if (response.statusCode === 302 || response.statusCode === 301) {
      https.get(response.headers.location, (redirectResponse) => {
        redirectResponse.pipe(file);
        file.on('finish', () => {
          file.close();
          fs.chmodSync(binaryPath, '755');
          console.log('✓ Binary downloaded successfully');
        });
      });
    } else if (response.statusCode === 200) {
      response.pipe(file);
      file.on('finish', () => {
        file.close();
        fs.chmodSync(binaryPath, '755');
        console.log('✓ Binary downloaded successfully');
      });
    } else {
      fs.unlinkSync(binaryPath);
      console.error(`Failed to download: HTTP ${response.statusCode}`);
      process.exit(1);
    }
  }).on('error', (err) => {
    if (fs.existsSync(binaryPath)) fs.unlinkSync(binaryPath);
    console.error('Download failed:', err.message);
    process.exit(1);
  });
}

try {
  downloadBinary();
} catch (err) {
  console.error('Installation failed:', err.message);
  process.exit(1);
}
