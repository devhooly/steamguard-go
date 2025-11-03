# steamguard-cli (Go version)

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)

A command-line utility for setting up and using Steam Mobile Authenticator (Steam 2FA). Ported to Go from the [original Rust project](https://github.com/dyc3/steamguard-cli).

## ⚠️ Disclaimer

**This utility is in development. Use at your own risk!**

- ✅ Regularly back up your `maFiles` folder
- ✅ Write down your revocation code
- ⚠️ If you lose both maFiles and the revocation code, we can't help - your only option is Steam support

## 🚀 Features

- ✅ **Generate 2FA codes** - Quick generation of Steam Guard codes
- ✅ **Manage confirmations** - Accept/reject trades and market listings
- ✅ **Encrypted storage** - Protect your 2FA secrets with AES-256 encryption
- ✅ **Generate QR codes** - Export to other applications (KeeWeb, 1Password, Bitwarden)
- ✅ **SDA compatible** - Reads maFiles format from Steam Desktop Authenticator
- ✅ **Cross-platform** - Works on Linux, Windows, macOS

## 📦 Installation

### From source (requires Go 1.21+)

```bash
git clone https://github.com/devhooly/steamguard-go
cd steamguard-go
go build -o steamguard
sudo mv steamguard /usr/local/bin/  # Linux/macOS
```

### Binary releases

Download pre-built binaries from the [Releases](https://github.com/devhooly/steamguard-fork/releases) section.

## 📖 Usage

### maFiles location

`steamguard-cli` looks for the `maFiles/` folder in these locations:

**Linux:**

- `~/.config/steamguard-cli/maFiles/`
- `~/maFiles/`

**Windows:**

- `%APPDATA%\steamguard-cli\maFiles\`
- `%USERPROFILE%\maFiles\`

### Basic commands

#### Generate code (for the first account)

```bash
steamguard
```

#### Generate code for a specific account

```bash
steamguard -u username
```

#### View QR code for importing into other applications

```bash
steamguard qr               # First account
steamguard -u username qr   # Specific account
```

**Do not use:** Google Authenticator, Authy (they generate incorrect codes!)  
**Recommended:** KeeWeb, 1Password, Bitwarden

#### Manage trade confirmations

```bash
steamguard trade               # Show list
steamguard trade --accept      # Accept all
steamguard trade --reject      # Reject all
```

#### List all accounts

```bash
steamguard list
```

### Copy code to clipboard

**Linux:**

```bash
steamguard | xclip -selection clipboard
# or
steamguard | xsel --clipboard
```

**macOS:**

```bash
steamguard | pbcopy
```

**Windows (PowerShell):**

```powershell
steamguard | Set-Clipboard
```

## 🔒 Security

### Encryption

To enable encryption in `manifest.json`, set:

```json
{
  "encrypted": true,
  ...
}
```

When adding the first account, you'll be prompted for a password. Uses:

- **AES-256-CBC** for encryption
- **PBKDF2** (50,000 iterations) for key derivation from password
- Unique IV and salt for each account

### Memory clearing

All critical data (secrets, passwords) is wiped from memory after use.

## 📁 Project structure

```
steamguard-fork/
├── cmd/              # CLI commands
│   ├── root.go       # Main command and code generation
│   ├── setup.go      # New account setup
│   ├── qr.go         # QR code generation
│   ├── trade.go      # Confirmation management
│   └── list.go       # List accounts
├── internal/
│   ├── config/       # Configuration and paths
│   ├── crypto/       # Encryption/decryption
│   ├── manifest/     # Work with maFiles
│   ├── qrcode/       # QR code generation
│   ├── steamapi/     # Steam API client
│   └── steamguard/   # TOTP generator for Steam
├── main.go           # Entry point
├── go.mod
└── README.md
```

## 🛠️ Development

This project uses [Task](https://taskfile.dev/) for build automation.

### Installing Task

```bash
# Linux/macOS
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

# macOS (Homebrew)
brew install go-task/tap/go-task

# Windows (Chocolatey)
choco install go-task

# Go install
go install github.com/go-task/task/v3/cmd/task@latest
```

### Available commands

List all available commands:

```bash
task --list-all
```

Main commands:

```bash
task build       # Build binary
task build-all   # Build for all platforms
task test        # Run tests
task run         # Build and run
task clean       # Clean built files
task install     # Install to system
task fmt         # Format code
task lint        # Lint code
task dev         # Development mode with auto-reload
```

### Manual build

If you prefer not to use Task:

```bash
# Install dependencies
go mod download

# Build
go build -o steamguard

# Run
./steamguard --help
```

## 📋 Compatibility

Fully compatible with the `maFiles` format from [Steam Desktop Authenticator](https://github.com/Jessecar96/SteamDesktopAuthenticator). You can use existing maFiles without modifications.

## 🤝 Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss.

## 📄 License

- `steamguard-cli` (program) - **GPL-3.0**
- `internal/*` libraries - **MIT or Apache 2.0** (your choice)

## 🙏 Credits

- [dyc3/steamguard-cli](https://github.com/dyc3/steamguard-cli) - original Rust implementation
- [Jessecar96/SteamDesktopAuthenticator](https://github.com/Jessecar96/SteamDesktopAuthenticator) - maFiles format

## ⚡ Usage examples

### Create alias for quick access

```bash
# Linux/macOS
echo "alias sg='steamguard'" >> ~/.bashrc
source ~/.bashrc

# Now you can simply:
sg
```

### Use in scripts

```bash
#!/bin/bash
# auto-login.sh - automatic login with 2FA code

USERNAME="your_username"
PASSWORD="your_password"
CODE=$(steamguard -u $USERNAME)

echo "2FA code: $CODE"
# Your login logic here...
```

### Monitor confirmations

```bash
# Check every 5 minutes
watch -n 300 'steamguard trade'
```

## 🐛 Known issues

1. **Login and Setup functions** - require full Steam Auth protocol implementation (RSA, captcha, email codes)
2. **HTML parsing** - GetConfirmations requires parsing HTML confirmation page
