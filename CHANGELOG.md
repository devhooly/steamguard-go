# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Full TOTP generator implementation for Steam Guard
- Support for reading/writing maFiles format (compatible with SDA)
- Data encryption using AES-256-CBC
- QR code generation for exporting 2FA secrets
- CLI commands: code generation, confirmation management, account list
- Support for Linux, Windows, macOS
- English documentation
- Usage examples and integrations
- GitHub Actions for automated builds and releases
- Makefile for simplified development

### In development

- Full Steam Auth protocol implementation for setup command
- HTML confirmation page parsing
- System keyring integration for secure password storage
- Automated tests
- QR code login support
- Backup creation command

## [0.1.0] - 2025-10-28

### Added

- Initial project version
- Port of main functionality from Rust to Go
- Basic project structure
- CLI interface based on Cobra
- Compatibility with maFiles format

[Unreleased]: https://github.com/devhooly/steamguard-fork/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/devhooly/steamguard-fork/releases/tag/v0.1.0
