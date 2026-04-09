# envoy-cli

A lightweight CLI for managing and syncing .env files across local and remote environments with encryption support.

## Features

- 🔒 AES-256 encryption for sensitive environment variables
- 🔄 Sync .env files across multiple environments (dev, staging, prod)
- 📦 Store encrypted configs in Git safely
- 🚀 Simple CLI interface with minimal dependencies
- 🔑 Support for multiple encryption keys per environment

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or download pre-built binaries from the [releases page](https://github.com/yourusername/envoy-cli/releases).

## Quick Start

```bash
# Initialize a new envoy config
envoy init

# Encrypt your .env file
envoy encrypt .env --output .env.encrypted

# Decrypt to a specific environment
envoy decrypt .env.encrypted --env production

# Sync to remote storage (S3, GCS, or custom backends)
envoy sync push --env production

# Pull latest config from remote
envoy sync pull --env staging
```

## Configuration

Create an `envoy.yml` file in your project root:

```yaml
environments:
  - name: production
    key_id: prod-key-2024
  - name: staging
    key_id: staging-key-2024

storage:
  type: s3
  bucket: my-envoy-configs
```

## License

MIT License - see [LICENSE](LICENSE) for details.