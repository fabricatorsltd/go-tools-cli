# go-tools-cli

Interactive project bootstrapper for the [go-tools](https://github.com/mirkobrombin/go-tools) ecosystem by [Fabricators](https://github.com/fabricatorsltd).

Scaffolds new Go projects in seconds — picks modules, generates `go.mod`, boilerplate stubs, and wired examples — all driven by a terminal wizard or fully from flags.

```
◆ go-tools-cli
  Bootstrap a new Go project with the go-tools ecosystem

  Step 1/7

  Project name:
  > my-api█

  Press Enter to continue • Esc to go back
```

---

## Installation

### via `go install` (recommended)

```bash
go install github.com/fabricatorsltd/go-tools-cli@latest
```

### Pre-built binaries

Download from [Releases](https://github.com/fabricatorsltd/go-tools-cli/releases):

| Platform      | Binary                           |
| ------------- | -------------------------------- |
| Linux amd64   | `go-tools-cli-linux-amd64`       |
| Linux arm64   | `go-tools-cli-linux-arm64`       |
| Windows amd64 | `go-tools-cli-windows-amd64.exe` |
| Windows arm64 | `go-tools-cli-windows-arm64.exe` |
| macOS amd64   | `go-tools-cli-darwin-amd64`      |
| macOS arm64   | `go-tools-cli-darwin-arm64`      |

A `checksums.txt` (SHA-256) is attached to every release.

---

## Usage

### Interactive wizard

```bash
go-tools-cli init
```

Launches a 7-step terminal wizard:

| Step | Prompt         | Description                                    |
| ---- | -------------- | ---------------------------------------------- |
| 1    | Project name   | Name of your Go project                        |
| 2    | Module path    | Go module path (e.g. `github.com/acme/my-app`) |
| 3    | Preset         | Starting template (see below)                  |
| 4    | Modules        | Add / remove modules from the preset           |
| 5    | Module options | Per-module choices (e.g. broker for go-relay)  |
| 6    | Output depth   | How much code to generate                      |
| 7    | Confirm        | Review and confirm before writing files        |

**Keyboard shortcuts:** `↑ ↓` navigate · `Space` toggle · `Enter` confirm · `Esc` go back · `Ctrl+C` abort

### Skip steps with flags

Any flag pre-fills a step and skips it in the wizard:

```bash
go-tools-cli init [flags] [-o <output-dir>]
```

| Flag       | Short | Description                      |
| ---------- | ----- | -------------------------------- |
| `--name`   | `-n`  | Project name                     |
| `--module` | `-m`  | Go module path                   |
| `--preset` | `-p`  | Preset name (see aliases below)  |
| `--depth`  | `-d`  | Output depth (see aliases below) |
| `--output` | `-o`  | Output directory (default: `.`)  |

---

## Presets

| Preset             | Alias       | Included modules                                                    |
| ------------------ | ----------- | ------------------------------------------------------------------- |
| **API server**     | `api`       | foundation, logger, conf-builder, module-router, auth, guard, httpx |
| **Worker service** | `worker`    | foundation, logger, conf-builder, relay, worker, signal, retry      |
| **Full-stack**     | `fullstack` | All 22 go-tools modules                                             |
| **Minimal**        | `minimal`   | foundation, logger, conf-builder                                    |

---

## Output depths

| Depth            | Alias         | What gets generated                                  |
| ---------------- | ------------- | ---------------------------------------------------- |
| **Minimal**      | `minimal`     | `go.mod` + `main.go` stub + `README.md`              |
| **Boilerplate**  | `boilerplate` | Minimal + `internal/` stubs for each selected module |
| **Full example** | `full`        | Boilerplate + wired `cmd/` entrypoint + `Makefile`   |

---

## Examples

### Fully interactive

```bash
go-tools-cli init -o ./my-project
```

### API server, boilerplate, choose modules interactively

```bash
go-tools-cli init -n my-api -m github.com/acme/my-api -p api -d boilerplate -o ./my-api
```

### Worker service, skip everything, generate full example

```bash
go-tools-cli init -n my-worker -m github.com/acme/my-worker -p worker -d full -o ./my-worker
```

### Minimal scaffold, just go.mod + main.go

```bash
go-tools-cli init -n my-lib -m github.com/acme/my-lib -p minimal -d minimal -o ./my-lib
```

After scaffolding:

```bash
cd my-api
go mod tidy
go run .
```

---

## Available modules

| Module                | Category      | Description                                              |
| --------------------- | ------------- | -------------------------------------------------------- |
| `go-foundation`       | core          | Zero-dependency primitives (DI, resiliency, safemap…)    |
| `go-logger`           | observability | Structured logging with pluggable sinks                  |
| `go-metrics`          | observability | Counter abstraction (Prometheus-compatible)              |
| `go-cli-builder/v2`   | cli           | Declarative struct-tag-driven CLI builder                |
| `go-conf-builder/v2`  | cli           | Multi-source config loader (env, file, flags)            |
| `go-struct-flags/v2`  | cli           | Flag-to-struct binding via struct tags                   |
| `go-module-router/v2` | http          | Transport-agnostic router with DI                        |
| `go-httpx`            | http          | Middleware-enhanced HTTP client                          |
| `go-auth`             | auth          | HMAC-SHA256 token signing and verification               |
| `go-guard`            | auth          | Declarative attribute-based access control (ABAC)        |
| `go-secrets`          | auth          | Secret store abstraction (env, vault, file)              |
| `go-signal/v2`        | async         | Type-safe in-process event bus                           |
| `go-relay/v2`         | async         | Durable async job processing (Memory/Redis/NATS/Warp)    |
| `go-worker`           | async         | Fixed-size worker pool with graceful shutdown            |
| `go-retry`            | async         | Compact retry helper with exponential backoff            |
| `go-revert/v2`        | async         | Saga / compensation workflows with panic-safe rollback   |
| `go-state-flow`       | state         | Declarative finite-state machine via struct tags         |
| `go-lock`             | infra         | Distributed locking (Redis, Redlock, etcd)               |
| `go-warp`             | infra         | L1/L2 cache with distributed sync bus                    |
| `go-slipstream`       | infra         | Embedded Bitcask+Raft database                           |
| `go-plugin`           | plugin        | Plugin registry with lifecycle management                |
| `go-wormhole`         | data          | EF-style ORM with code-first migrations and Unit of Work |

> **Tip:** Use the [go-tools MCP server](https://go-tools.fabricators.ltd/mcp) to query module documentation directly from your AI assistant.

---

## Contributing

This tool is part of the [go-tools](https://github.com/mirkobrombin/go-tools) monorepo.
Issues and PRs welcome at [fabricatorsltd/go-tools-cli](https://github.com/fabricatorsltd/go-tools-cli).
