# OpsOrch {ADAPTER_NAME} Adapter

<!-- 
TEMPLATE INSTRUCTIONS:
This is the OpsOrch Adapter Starter template. When creating a new adapter:
1. Replace all {PLACEHOLDER} values with your adapter-specific information
2. Remove sections marked [Optional] if not applicable
3. Keep this comment block during development, remove before publishing
4. Follow the standardized structure to ensure consistency across all adapters
-->

[![Version](https://img.shields.io/github/v/release/opsorch/{REPO_NAME})](https://github.com/opsorch/{REPO_NAME}/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/opsorch/{REPO_NAME})](https://github.com/opsorch/{REPO_NAME}/blob/main/go.mod)
[![License](https://img.shields.io/github/license/opsorch/{REPO_NAME})](https://github.com/opsorch/{REPO_NAME}/blob/main/LICENSE)
[![CI](https://github.com/opsorch/{REPO_NAME}/workflows/CI/badge.svg)](https://github.com/opsorch/{REPO_NAME}/actions)

<!-- Replace with one-sentence description of what your adapter does -->
This adapter integrates OpsOrch with {EXTERNAL_SYSTEM}, enabling {brief description of capabilities}.

## Quick Start

This is a copy/paste-friendly starter for building an OpsOrch external adapter or plugin. Follow these steps to create your own adapter:

1. **Clone/fork this repository**
   ```bash
   git clone https://github.com/opsorch/opsorch-adapter.git opsorch-{yourprovider}-adapter
   cd opsorch-{yourprovider}-adapter
   ```

2. **Rename the module**
   - Update `module` in `go.mod` to your repo path (e.g., `github.com/yourorg/opsorch-yourprovider-adapter`)
   - Keep the `replace` directive pointing at your local `opsorch-core` during development
   - Remove the `replace` directive when publishing

3. **Choose your provider name**
   - Pick a provider name (e.g., `pagerduty`, `datadog`, `mytool`)
   - Update references to `example` in code to match your provider name

4. **Implement your capability**
   - Implement the capability interfaces you need (incident/log/metric/ticket/messaging/alert/service/secret)
   - An incident example is provided as a starting point
   - Register your provider in `init()` with the registry for in-process use

5. **Update this README**
   - Replace all {PLACEHOLDER} values with your adapter-specific information
   - Follow the sections below to document your adapter
   - Remove template comments before publishing

## Capabilities

<!-- List all capabilities your adapter implements -->

This adapter provides the following capabilities:

1. **{Capability 1}**: {Brief description}
<!-- Add more capabilities as needed -->

## Features

<!-- Organize features by capability if multi-capability adapter -->

- **{Feature 1}**: {Description}
- **{Feature 2}**: {Description}
- **{Feature 3}**: {Description}

### Version Compatibility

- **Adapter Version**: {version}
- **Requires OpsOrch Core**: >={version}
- **{External System} Version**: {supported versions}
- **Go Version**: {version}+

## Configuration

<!-- Document all configuration fields your adapter requires -->

The {capability} adapter requires the following configuration:

| Field | Type | Required | Description | Default |
|-------|------|----------|-------------|---------|
| `{fieldName}` | {type} | Yes/No | {Description} | `{default}` |
<!-- Add all configuration fields -->

### Authentication Setup

<!-- Provide step-by-step instructions for obtaining credentials -->

#### 1. {First Step Title}

{Detailed instructions for obtaining credentials}

#### 2. {Second Step Title}

{Additional setup steps}

### Example Configuration

**JSON format:**
```json
{
  "{fieldName}": "{example-value}",
  "{fieldName}": "{example-value}"
}
```

**Environment variables:**
```bash
export OPSORCH_{CAPABILITY}_PLUGIN=/path/to/bin/{capability}plugin
export OPSORCH_{CAPABILITY}_CONFIG='{"fieldName":"value"}'
```

## Field Mapping

<!-- This section applies to data adapters: log, metric, alert, incident, service -->
<!-- Skip this section for messaging, ticket, and secret adapters unless they transform data -->

### Query Mapping

How OpsOrch query fields map to {External System} API parameters:

| OpsOrch Field | {External System} Parameter | Transformation | Notes |
|---------------|----------------------------|----------------|-------|
| `{field}` | `{parameter}` | {How value is transformed} | {Additional context} |

### Response Normalization

How {External System} response fields map to OpsOrch schema:

| {External System} Field | OpsOrch Field | Transformation | Notes |
|------------------------|---------------|----------------|-------|
| `{field}` | `{field}` | {How value is transformed} | {Additional context} |

### Metadata Fields

Provider-specific fields stored in the `metadata` map:

| Metadata Key | Source Field | Type | Description |
|--------------|--------------|------|-------------|
| `{key}` | `{source}` | {type} | {What this field contains} |

### Known Limitations

<!-- Explicitly document API limitations, unsupported features, workarounds -->

1. **{Limitation Title}**: {Description of limitation and impact}
2. **{Limitation Title}**: {Description of limitation and impact}

## Usage

### In-Process Mode

Import the adapter for side effects to register it with OpsOrch Core:

```go
import (
    _ "github.com/opsorch/{REPO_NAME}/{capability}"
)
```

Configure via environment variables:

```bash
export OPSORCH_{CAPABILITY}_PROVIDER={provider-name}
export OPSORCH_{CAPABILITY}_CONFIG='{"fieldName":"value"}'
```

### Plugin Mode

Build the plugin binary:

```bash
make plugin
```

Configure OpsOrch Core to use the plugin:

```bash
export OPSORCH_{CAPABILITY}_PLUGIN=/path/to/bin/{capability}plugin
export OPSORCH_{CAPABILITY}_CONFIG='{"fieldName":"value"}'
```

### Docker Deployment

Download pre-built plugin binaries from [GitHub Releases](https://github.com/opsorch/{REPO_NAME}/releases):

```dockerfile
FROM ghcr.io/opsorch/opsorch-core:latest
WORKDIR /opt/opsorch

# Download plugin binary
ADD https://github.com/opsorch/{REPO_NAME}/releases/download/v{VERSION}/{capability}plugin-linux-amd64 ./plugins/{capability}plugin
RUN chmod +x ./plugins/{capability}plugin

# Configure plugin
ENV OPSORCH_{CAPABILITY}_PLUGIN=/opt/opsorch/plugins/{capability}plugin
```

## Development

### Prerequisites

- Go {version}+ (e.g., 1.21+)
- {External system access requirements}
- {Any other prerequisites}

### Building

```bash
# Download dependencies
go mod download

# Run unit tests
make test

# Build all packages
make build

# Build plugin binary
make plugin

# Run integration tests (requires credentials)
make integ
```

### Testing

**Unit Tests:**
```bash
make test
```

**Integration Tests:**

Integration tests run against a real {External System} instance.

**Prerequisites:**
- {Prerequisite 1}
- {Prerequisite 2}

**Setup:**
```bash
# Set required environment variables
export {ENV_VAR_1}="value"
export {ENV_VAR_2}="value"

# Run all integration tests
make integ

# Or run specific capability tests
make integ-{capability}
```

**What the tests do:**
- {Description of test 1}
- {Description of test 2}

**Expected behavior:**
- {Expected outcome 1}
- {Expected outcome 2}

### Project Structure

```
{REPO_NAME}/
├── {capability}/              # {Capability} provider implementation
│   ├── {provider}.go         # Core provider logic
│   └── {provider}_test.go    # Unit tests
├── cmd/
│   └── {capability}plugin/   # Plugin entrypoint
│       └── main.go
├── integ/                     # Integration tests
│   └── {capability}.go
├── version.go                 # Adapter version metadata
├── Makefile                   # Build targets
├── go.mod
└── README.md
```

**Key Components:**

- **{capability}/{provider}.go**: Implements {capability}.Provider interface, handles {description of what it does}
- **cmd/{capability}plugin**: JSON-RPC plugin wrapper for {capability} provider
- **integ/{capability}.go**: End-to-end integration tests against live {External System} instance
- **version.go**: Adapter version metadata; keep in sync with releases

## CI/CD & Pre-Built Binaries

The repository includes GitHub Actions workflows:

- **CI** (`ci.yml`): Runs tests and linting on every push/PR to main
- **Release** (`release.yml`): Manual workflow that:
  - Runs tests and linting
  - Creates version tags (patch/minor/major)
  - Builds multi-arch binaries (linux-amd64, linux-arm64, darwin-amd64, darwin-arm64)
  - Publishes binaries as GitHub release assets

### Downloading Pre-Built Binaries

Pre-built plugin binaries are available from [GitHub Releases](https://github.com/opsorch/{REPO_NAME}/releases).

**Supported platforms:**
- Linux (amd64, arm64)
- macOS (amd64, arm64)

**Example usage in Dockerfile:**
```dockerfile
FROM ghcr.io/opsorch/opsorch-core:latest
WORKDIR /opt/opsorch

ADD https://github.com/opsorch/{REPO_NAME}/releases/download/v{VERSION}/{capability}plugin-linux-amd64 ./plugins/{capability}plugin
RUN chmod +x ./plugins/{capability}plugin

ENV OPSORCH_{CAPABILITY}_PLUGIN=/opt/opsorch/plugins/{capability}plugin
```

## Plugin RPC Contract

OpsOrch Core communicates with the plugin over stdin/stdout using JSON-RPC.

### Message Format

**Request:**
```json
{
  "method": "{capability}.{operation}",
  "config": { /* decrypted configuration */ },
  "payload": { /* method-specific request body */ }
}
```

**Response:**
```json
{
  "result": { /* method-specific result */ },
  "error": "optional error message"
}
```

### Configuration Injection

The `config` field contains the decrypted configuration map from `OPSORCH_{CAPABILITY}_CONFIG`. The plugin receives this on every request, so it never stores secrets on disk.

### Supported Methods

<!-- Document each RPC method with request/response examples -->

#### {capability}.{operation}

{Description of what this method does}

**Request:**
```json
{
  "method": "{capability}.{operation}",
  "config": { /* config fields */ },
  "payload": {
    "{field}": "{value}"
  }
}
```

**Response:**
```json
{
  "result": {
    "{field}": "{value}"
  }
}
```

<!-- Repeat for each supported method -->

## Security Considerations

1. **Never log credentials**: Avoid logging the config or sensitive fields in application logs
2. **Rotate credentials regularly**: Follow your organization's security policy for credential rotation
3. **Use environment variables**: Store config in secure environment variables or secrets management systems
4. **Restrict file permissions**: If storing config in files, use appropriate permissions (e.g., 0600)
5. **Use API tokens**: Prefer API tokens over passwords where supported by the external system
6. **Validate TLS certificates**: Ensure TLS certificate validation is enabled for production use

## Production Guidance

When building production adapters, follow these best practices:

- **Keep adapters stateless**: Never persist tokens or credentials on disk
- **Normalize responses**: Map external system responses to the current `opsorch-core/schema` types
- **Use metadata for extras**: Store provider-specific fields in the `metadata` map
- **Wrap provider errors**: Avoid leaking raw API responses or secrets in error messages
- **Version your adapter**: Keep `AdapterVersion` and `RequiresCore` in sync with releases
- **Build static binaries**: For plugins, build static binaries and checksum them
- **Pin dependencies**: Pin OpsOrch Core compatibility in your `go.mod`

## License

{License type - typically Apache 2.0 or MIT}

See LICENSE file for details.
