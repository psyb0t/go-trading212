package trading212_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-trading212"
)

// This is the README usage example. It talks to the live API, so it has no
// Output comment and only has to compile.
//
//nolint:testableexamples // needs real credentials, compile check only
func Example() {
	auth := trading212.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.SetBasicAuth(
				os.Getenv("TRADING212_API_KEY"),
				os.Getenv("TRADING212_API_SECRET"),
			)

			return nil
		},
	)

	client, err := trading212.NewClientWithResponses(
		"https://demo.trading212.com",
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.GetAccountSummaryWithResponse(context.Background())
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
