# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

Gravity is a fully-replicated DNS, DHCP, and TFTP server backed by etcd. The server is a single Go binary (`gravity`) plus a Lit/TypeScript web UI. All persistent state lives in etcd; cluster nodes are peers and replicate through embedded etcd.

## Common commands

All targets are in the top-level `Makefile`. `.SHELLFLAGS += -x -e -o pipefail` is set, so Make echoes every command.

### Server / Go

- `make run` — run the dev server (sets `DEBUG=true`, `LISTEN_ONLY=true`, binds `0.0.0.0`). Requires `internal/resources/{macoui,blocky,tftp}/` — the targets that fetch them require `tshark` and network access; they're skipped on subsequent runs because they exist.
- `make test-env-start` / `make test-env-stop` — bring up the etcd (port 2385) + dex + coredns + nginx stack used by the unit tests (`hack/tests/docker-compose.yml`).
- `make test` — runs unit tests. Requires `test-env-start` first. Tests run with `-p 1` (no parallel packages) and exclude generated packages (`api/`, `cmd/`, `pkg/externaldns/generated/`). Coverage HTML is written to `coverage.html`.
- `make test-local` — sets `TEST_COUNT=100` + `-shuffle=on -failfast`; chain it with `test` (`make test-local test`) for flake hunting.
- Single test: `go test -v -run TestName ./pkg/roles/dns/...` — must export the same env as `make test` (`BOOTSTRAP_ROLES`, `ETCD_ENDPOINT=localhost:2385`, `DEBUG=true`, `LISTEN_ONLY=true`, `CI=true`).
- `make test-e2e` — builds a `gravity:e2e-test` Docker image and runs tests under `tests/` with build tag `e2e` via testcontainers. Slow.
- `make bench` — benchmark suite.
- `make lint` — `golangci-lint run` (5000s timeout) + `web-lint`.
- `make bin/gravity-cli` — build the standalone CLI binary at `bin/gravity-cli`.

### Web

- `make web-install` — `npm ci` inside `web/`.
- `make web-watch` — rollup watch mode.
- `make web-build` — production build into `web/dist/` (embedded into the Go binary via `web/static.go`).
- `make web-lint` — runs prettier, eslint, tsc, and lit-analyzer; CI requires zero warnings (`--max-warnings 0`).

### Code generation

The OpenAPI schema (`schema.yml`) is the source of truth for the API clients. Regenerate after changing API handlers:

- `make gen-build` — runs `gravity generateSchema schema.yml` to refresh the schema from the running code, then `git add`s it.
- `make gen-client-go` — regenerates `api/` (Go client) via openapi-generator in Docker.
- `make gen-client-ts` — regenerates `gen-ts-api/` (TypeScript client), then copies into `web/node_modules/gravity-api`.
- `make gen-external-dns` — regenerates the external-dns webhook server stubs in `pkg/externaldns/generated/`.
- `make gen-proto` — regenerates protobuf bindings under `protobuf/`.

The release flow (`make release`) chains `gen-build → gen-clean → gen-client-go → gen-external-dns → gen-go-tidy → gen-client-ts-publish → gen-tag`.

## Architecture

### Instance + roles

Gravity is structured as an `Instance` (`pkg/instance/`) that hosts a set of `Role`s (`pkg/roles/`). Each role registers itself in `init()` via `roles.Register(name, constructor)` and implements:

```go
type Role interface {
    Start(ctx context.Context, config []byte) error
    Stop()
}
```

The roles in this repo: `etcd`, `api`, `dns`, `dhcp`, `tftp`, `discovery`, `backup`, `monitoring`, `tsdb`, `debug`. Which roles a node runs is driven by the env var `BOOTSTRAP_ROLES` (default includes everything) and overridden per-instance by a value stored in etcd at `/<prefix>/instance/<id>/roles`.

`Instance.bootstrap()` (`pkg/instance/instance.go`) constructs each role with `roles.GetRole(id)(roleInstance)`, watches the role's config key in etcd, and on every change dispatches an `EventTopicRoleRestart` event that stops + restarts the role. Roles never read config from disk at runtime — they receive raw bytes (typically JSON) on `Start`. If a role implements `MigratableRole`, its `RegisterMigrations()` runs before `Start`, hooked via `RoleMigrator`.

The `etcd` role is special: if `etcd` is in `BOOTSTRAP_ROLES`, the embedded etcd starts first (`startEtcd`) before any other role bootstraps, since they all need a KV client.

### Events

`Instance` exposes an in-process event bus (`AddEventListener` / `DispatchEvent`) used to wire roles to each other (e.g. DNS subscribes to `EventTopicDHCPLeasePut` to auto-create A records, see `pkg/roles/dns/dhcp_event_handler.go`). Topics are constants in each role's `types/` package. Stay loosely coupled — don't import another role's package to call into it; emit/consume events instead.

### Storage

`pkg/storage` wraps the etcd client. Always build keys via `client.Key("part", "part")` rather than concatenating strings — the wrapper applies a namespace prefix (`ETCD_PREFIX`, default `/gravity`) and supports `Prefix(true)` for range gets. `pkg/storage/watcher` is a generic typed wrapper that decodes KV events into your value type and keeps an in-memory mirror — used for things like DNS zones (`pkg/roles/dns/role.go`). Migrations hook `Get`/`Put`/`Delete` via `StorageHook`.

### API + web

The `api` role uses [`swaggest/rest`](https://github.com/swaggest/rest) to register handlers and emit OpenAPI. `cmd/server/generateSchema.go` boots the server in schema-emit mode and dumps `schema.yml`. From there, `make gen-client-go` and `make gen-client-ts` produce the Go (`api/`) and TypeScript (`gen-ts-api/`, copied into `web/node_modules/gravity-api`) clients. Treat `api/`, `gen-ts-api/`, and `pkg/externaldns/generated/` as generated — never hand-edit; rerun the relevant `make gen-*` target.

The web UI is Lit-based custom elements (`web/src/`), bundled with rollup, and served from the Go binary via `web/static.go` (the API role mounts it). Page routes live under `web/src/pages/`.

### CLI

`cmd/cli/` adds subcommands to the same `gravity` binary — they talk to a running server over its REST API (default Unix socket `/var/run/gravity.sock`, or the cwd's `gravity.sock` in `DEBUG=true`). Useful subcommands: `tokens`, `users`, `export`, `import`, `convert` (BIND/MS-DHCP/OPNsense), `etcdctl`, `snapshot`, `health`, `debug`. The standalone `bin/gravity-cli` is just the same code built from `cmd/cli/main`.

### Configuration

All runtime config is environment variables, parsed via `pkg/extconfig` (using `Netflix/go-env`). `extconfig.Get()` is a process-global singleton — call it from anywhere instead of plumbing config through. Notable vars: `INSTANCE_IDENTIFIER`, `INSTANCE_IP`, `INSTANCE_LISTEN`, `ETCD_PREFIX`, `ETCD_ENDPOINT`, `ETCD_JOIN_CLUSTER`, `BOOTSTRAP_ROLES`, `DATA_PATH`, `DEBUG`, `LISTEN_ONLY` (skips binding privileged ports), `CI`.

## Testing notes

- Unit tests assume an etcd reachable at `localhost:2385` (provided by `make test-env-start`). They run with `-p 1` because they share that etcd and each test typically wipes the namespace.
- E2E tests are gated behind `//go:build e2e` and use testcontainers — they boot a real `gravity:e2e-test` Docker image. Don't add e2e-only deps to non-tagged files.
- The `web/static.go` build embeds `web/dist/`. If you're working on Go-only changes, you don't need to rebuild the web UI; the dev `make run` target doesn't depend on it either.
