# OpsOrch Adapter Starter

This repo is a copy/paste-friendly starter for building an OpsOrch external adapter or plugin. Fork it, rename the module, and implement your provider logic.

## Quick start

1) Clone/fork and rename the module:
   - Update `module` in `go.mod` to your repo path.
   - Keep the `replace` pointing at your local `opsorch-core` during development; remove it when publishing.
2) Pick a provider name (e.g., `pagerduty`, `mytool`) and update references to `example` in code to match.
3) Implement the capability surfaces you need (incident/log/metric/...). An incident example is provided.
4) Register your provider in `init()` with the registry for in-process use.
5) If you need a plugin binary, build `cmd/incidentplugin` (or your renamed path) and point OpsOrch Core at it with `OPSORCH_<CAP>_PLUGIN`.

## Repository layout

- `incident/example_provider.go`: Minimal incident provider showing config parsing, normalization, and timeline support.
- `cmd/incidentplugin/main.go`: JSON-RPC plugin entrypoint that wraps the example provider for `incident.*` methods.
- `version.go`: Adapter version metadata; keep in sync with releases.
- `Makefile`: Convenience targets for build/test/plugin builds.

## Building

```bash
make test    # run unit tests
make build   # build all packages
make plugin  # builds ./bin/incidentplugin
make integ   # run integration tests
```

## Configuration

The sample incident provider expects a config like:

```json
{
  "source": "demo",
  "defaultSeverity": "sev3"
}
```

Fields are illustrative; change to whatever your provider needs. OpsOrch Core will hand you a decrypted `map[string]any` at runtime.

## Production guidance

- Keep adapters stateless; never persist tokens.
- Normalize responses to the current `opsorch-core/schema` types; use `metadata` for provider-specific extras.
- Wrap provider errors; avoid leaking raw API responses or secrets.
- Version your adapter (`AdapterVersion`, `RequiresCore`) and pin core compatibility.
- For plugins, build static binaries and checksum them; point `OPSORCH_<CAP>_PLUGIN` to the baked-in path in your container.
