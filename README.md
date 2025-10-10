# libp2p-keygen

A command-line tool for generating Ed25519 key pairs for libp2p with optional vanity peer ID support.

## Features

- 🔐 Generate Ed25519 private keys compatible with libp2p
- 🎨 Vanity peer ID generation - search for custom patterns
- ⚡ Fast key generation with progress tracking
- 🔧 Case-sensitive and case-insensitive pattern matching
- 📦 Cross-platform binaries (macOS, Linux, Windows)
- 🎯 Simple CLI interface

## Installation

### Download Pre-built Binaries (Recommended)

Download the latest release from [GitHub Releases](https://github.com/bsv-blockchain/libp2p-keygen/releases).

**Quick download for your platform:**

**macOS (Apple Silicon):**
```bash
curl -LO https://github.com/bsv-blockchain/libp2p-keygen/releases/latest/download/libp2p-keygen-darwin-arm64
chmod +x libp2p-keygen-darwin-arm64
sudo mv libp2p-keygen-darwin-arm64 /usr/local/bin/libp2p-keygen
```

**macOS (Intel):**
```bash
curl -LO https://github.com/bsv-blockchain/libp2p-keygen/releases/latest/download/libp2p-keygen-darwin-amd64
chmod +x libp2p-keygen-darwin-amd64
sudo mv libp2p-keygen-darwin-amd64 /usr/local/bin/libp2p-keygen
```

**Linux (amd64):**
```bash
curl -LO https://github.com/bsv-blockchain/libp2p-keygen/releases/latest/download/libp2p-keygen-linux-amd64
chmod +x libp2p-keygen-linux-amd64
sudo mv libp2p-keygen-linux-amd64 /usr/local/bin/libp2p-keygen
```

**Linux (arm64):**
```bash
curl -LO https://github.com/bsv-blockchain/libp2p-keygen/releases/latest/download/libp2p-keygen-linux-arm64
chmod +x libp2p-keygen-linux-arm64
sudo mv libp2p-keygen-linux-arm64 /usr/local/bin/libp2p-keygen
```

**Windows (amd64):**
Download `libp2p-keygen-windows-amd64.exe` from the [releases page](https://github.com/bsv-blockchain/libp2p-keygen/releases).

### Build From Source

```bash
# Clone the repository
git clone https://github.com/bsv-blockchain/libp2p-keygen.git
cd libp2p-keygen

# Build and install
make install
```

Or build for a specific platform:
```bash
# Build for current platform
make build

# Build for all platforms
make build-all
```

Binaries will be available in the `bin/` directory.

## Usage

### Basic Key Generation

Generate a single Ed25519 key pair:

```bash
libp2p-keygen
```

Output:
```
Private Key (hex):
8d89a7c2f5b4e1d9a3c6f7e8b2d1a9c4f5e6d7b8a9c0d1e2f3a4b5c6d7e8f9a0

Peer ID: 12D3KooWABCDefGHijKLmnoPQRsTUVwxYZ1234567890aBcDeFgHi

Connect to this peer:
  /ip4/0.0.0.0/tcp/9095/p2p/12D3KooWABCDefGHijKLmnoPQRsTUVwxYZ1234567890aBcDeFgHi
  /ip6/::/tcp/9095/p2p/12D3KooWABCDefGHijKLmnoPQRsTUVwxYZ1234567890aBcDeFgHi
```

### Vanity Peer ID Generation

Generate a key with a specific pattern in the peer ID:

```bash
# Case-insensitive search (default)
libp2p-keygen --vanity cool

# Case-sensitive search
libp2p-keygen --vanity Cool --case-sensitive

# With verbose output
libp2p-keygen --vanity bitcoin --verbose
```

Example output:
```
Searching for vanity peer ID containing 'cool'...
Attempts: 15000 (1500.32 keys/sec)

✓ Found match after 15234 attempts (10.15 seconds, 1500 keys/sec)

Private Key (hex):
a1b2c3d4e5f6...

Peer ID: 12D3KooWCooLxYz...

Connect to this peer:
  /ip4/0.0.0.0/tcp/9095/p2p/12D3KooWCooLxYz...
```

### CLI Options

```
Usage:
  libp2p-keygen [flags]

Flags:
  --vanity string        Search for a vanity peer ID containing this pattern
  --case-sensitive       Enable case-sensitive pattern matching (default: false)
  --verbose             Show verbose output during vanity generation
  --port int            Port for libp2p host display (default: 9095)
  --version             Show version information
```

## Examples

```bash
# Simple generation
libp2p-keygen

# Find a peer ID containing "bsv"
libp2p-keygen --vanity bsv

# Find a peer ID with exact case "BSV"
libp2p-keygen --vanity BSV --case-sensitive

# Generate with custom port display
libp2p-keygen --port 8080

# Show version
libp2p-keygen --version
```

## Understanding Vanity Generation

Vanity peer ID generation works by repeatedly generating keys until finding one whose peer ID contains your desired pattern.

**Performance expectations:**
- 2-3 character patterns: Usually found in seconds
- 4-5 character patterns: May take minutes
- 6+ character patterns: Can take hours or days
- Case-sensitive searches take longer than case-insensitive

**Tips:**
- Start with shorter patterns (3-4 characters)
- Use case-insensitive matching for faster results
- Use `--verbose` to see progress
- The search is probabilistic - time varies significantly

## Using Generated Keys

### Teranode

The tool outputs ready-to-use configuration for Teranode. Simply copy the private key and use it in one of these ways:

**Environment Variable:**
```bash
export p2p_private_key="708fc975296cbb7d6bd1b572bc018a7010deb08e9e931608e177d1480f85d563..."
```

**Configuration File (settings_local.conf):**
```conf
p2p_private_key = "708fc975296cbb7d6bd1b572bc018a7010deb08e9e931608e177d1480f85d563..."
```

**Important Notes:**
- The Peer ID is automatically derived from your private key
- You only need to store the private key - the Peer ID doesn't need to be stored separately
- Keep your private key secure and never share it

### Go (libp2p)

```go
import (
    "encoding/hex"
    "github.com/libp2p/go-libp2p"
    "github.com/libp2p/go-libp2p/core/crypto"
)

func createHost(hexKey string) (host.Host, error) {
    keyBytes, _ := hex.DecodeString(hexKey)
    privateKey, _ := crypto.UnmarshalEd25519PrivateKey(keyBytes)

    return libp2p.New(
        libp2p.Identity(privateKey),
    )
}
```

### Other libp2p Applications

Many libp2p applications accept hex-encoded private keys in configuration:

```conf
# settings.conf
p2p_private_key = "8d89a7c2f5b4e1d9a3c6f7e8b2d1a9c4..."
```

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Run Examples

```bash
make examples
```

### Build for All Platforms

```bash
make build-all
```

This creates binaries for:
- macOS (Intel and Apple Silicon)
- Linux (amd64 and arm64)
- Windows (amd64)

## Security Notes

- Generated keys use crypto/rand for secure randomness
- Private keys are displayed in hex format - keep them secure
- Never share your private key
- Consider using environment variables or secure key storage in production

## Creating a Release

To create a new release with binaries:

```bash
# Tag the version
git tag v1.0.0

# Push the tag
git push origin v1.0.0
```

GitHub Actions will automatically build binaries for all platforms and create the release.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

Copyright 2025 BSV Association

## Acknowledgments

- Built with [go-libp2p](https://github.com/libp2p/go-libp2p)
- Created for the BSV blockchain ecosystem
