package trading212_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-trading212"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIKey    = "test-key"
	testAPISecret = "test-secret"
)

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *trading212.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := trading212.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.SetBasicAuth(testAPIKey, testAPISecret)

			return nil
		},
	)

	client, err := trading212.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func assertBasicAuth(t *testing.T, r *http.Request) {
	t.Helper()

	user, pass, ok := r.BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, testAPIKey, user)
	assert.Equal(t, testAPISecret, pass)
}

func TestGetAccountSummary(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v0/equity/account/summary", r.URL.Path)
		assertBasicAuth(t, r)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"id":42,"currency":"EUR","totalValue":1234.56,` +
				`"cash":{"availableToTrade":100.25}}`,
		))
	})

	res, err := client.GetAccountSummaryWithResponse(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)

	assert.Equal(t, int64(42), *res.JSON200.Id)
	assert.Equal(t, "EUR", *res.JSON200.Currency)
	assert.InDelta(t, 1234.56, *res.JSON200.TotalValue, 0)
	assert.InDelta(t, 100.25, *res.JSON200.Cash.AvailableToTrade, 0)
}

func TestGetPositionsWithTicker(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v0/equity/positions", r.URL.Path)
		assert.Equal(t, "TICKER_X", r.URL.Query().Get("ticker"))
		assertBasicAuth(t, r)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`[{"quantity":2.5,"currentPrice":10.1,` +
				`"walletImpact":{"fxImpact":null}}]`,
		))
	})

	ticker := "TICKER_X"

	res, err := client.GetPositionsWithResponse(
		context.Background(),
		&trading212.GetPositionsParams{Ticker: &ticker},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	require.Len(t, *res.JSON200, 1)

	pos := (*res.JSON200)[0]
	assert.InDelta(t, 2.5, *pos.Quantity, 0)
	assert.InDelta(t, 10.1, *pos.CurrentPrice, 0)
	assert.Nil(t, pos.WalletImpact.FxImpact)
}

func TestPlaceMarketOrder(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v0/equity/orders/market", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assertBasicAuth(t, r)

		var body trading212.MarketRequest

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "TICKER_X", *body.Ticker)
		assert.InDelta(t, 3.0, *body.Quantity, 0)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":7,"status":"NEW","quantity":3}`))
	})

	ticker := "TICKER_X"
	quantity := 3.0

	res, err := client.PlaceMarketOrderWithResponse(
		context.Background(),
		trading212.MarketRequest{Ticker: &ticker, Quantity: &quantity},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)

	assert.Equal(t, int64(7), *res.JSON200.Id)
	assert.Equal(t, trading212.OrderStatus("NEW"), *res.JSON200.Status)
}
