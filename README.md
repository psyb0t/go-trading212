# go-trading212

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-trading212.svg)](https://pkg.go.dev/github.com/psyb0t/go-trading212)
[![CI](https://github.com/psyb0t/go-trading212/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-trading212/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-trading212/badges/version.svg)](https://github.com/psyb0t/go-trading212/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-trading212/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-trading212/badges/importers.svg)](https://github.com/psyb0t/go-trading212/blob/badges/importers.md)

Unofficial Go client for the [Trading 212 API](https://docs.trading212.com), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI document Trading 212 publishes. I am not affiliated with Trading 212.

The package is `trading212`, in the repository root. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-trading212
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-trading212"
)

func main() {
	auth := trading212.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.SetBasicAuth(os.Getenv("TRADING212_API_KEY"), os.Getenv("TRADING212_API_SECRET"))

		return nil
	})

	client, err := trading212.NewClientWithResponses("https://demo.trading212.com", auth)
	if err != nil {
		panic(err)
	}

	res, err := client.GetAccountSummaryWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code is in `example_test.go`, so `make test` checks that it compiles.

Method names come from the spec's operation IDs, so some are generic (`GetAll`, `Create`, `Update` and `Delete` are the pie endpoints) and a few carry a number suffix (`Orders1` is the order history, `PlaceStopOrder1` is the stop order, `PlaceStopOrder` is the stop-limit order).

The full API, with the spec's descriptions as doc comments, is in the [Go reference](https://pkg.go.dev/github.com/psyb0t/go-trading212) and in `go doc github.com/psyb0t/go-trading212` from inside your module.

## Auth

Trading 212 authenticates with HTTP Basic auth: the API key is the username and the API secret is the password. Both come from the API section of the account settings in the Trading 212 app. The spec also lists a legacy scheme where the bare API key goes in the `Authorization` header; to use it, set `req.Header.Set("Authorization", key)` in the request editor instead.

Keys are tied to one environment. Pick the base URL that matches the key:

| Environment | Base URL |
|---|---|
| Demo (paper trading) | `https://demo.trading212.com` |
| Live (real money) | `https://live.trading212.com` |

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads `https://docs.trading212.com/_bundle/api.yaml` to `api/openapi.yaml`.
- `scripts/fix-nullable-types.sh` writes `api/openapi.fixed.yaml`. The spec declares OpenAPI 3.0.1 but uses the 3.1 type list form for nullable fields (`type: [number, 'null']`), which oapi-codegen rejects. The script rewrites each one to `type: number` plus `nullable: true`.
- oapi-codegen runs with `config.yaml`. The config maps `number` to `float64`, because the spec gives prices and amounts no format and oapi-codegen would otherwise use `float32`.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release: a patch version for fixes and additions, a minor version for breaking changes.

## Agent integrations

The [skill](.agents/skills/go-trading212) teaches a coding agent how to use go-trading212 when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-trading212@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-trading212@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-trading212
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from Trading 212's published API specification. The specification itself is not redistributed here.
