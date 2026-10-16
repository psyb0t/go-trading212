// Package trading212 is a Go client for the Trading 212 Public API,
// generated with oapi-codegen from the OpenAPI document Trading 212 publishes.
package trading212

//go:generate bash scripts/fetch-spec.sh
//go:generate bash scripts/fix-nullable-types.sh
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.fixed.yaml
