# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.2, 2026-10-17

- Documentation only. The README and `AGENTS.md` describe the daily regeneration workflow: what it checks, how it fixes failing checks, how it picks the next version and which paths a fix may not touch.

## v0.1.1, 2026-10-16

- Added the `go-trading212` agent skill. The README's new Agent integrations section has the install commands for Claude Code, Codex and ClawHub.

## v0.1.0, 2026-10-16

- Added `trading212`, a REST client generated with oapi-codegen from Trading 212's OpenAPI document (Trading 212 Public API, spec version v0, OpenAPI 3.0.1). It covers all 17 paths: account summary, history and CSV exports, instrument metadata, orders, pies and positions.
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` downloads the current spec from docs.trading212.com, rewrites its 3.1 style nullable type lists into 3.0 form, and regenerates the client. The spec is not committed.
