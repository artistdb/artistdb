package e2e

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	httpClient := &http.Client{}

	// This should always be the first test in this suite.
	t.Run("health endpoint is reachable", func(t *testing.T) {
		do := func() bool {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/internal/health", nil)
			require.NoError(t, err)

			req.Close = true

			resp, err := httpClient.Do(req)
			if err != nil {
				return false
			}

			defer func() {
				require.NoError(t, resp.Body.Close())
			}()

			return resp.StatusCode == http.StatusOK
		}

		require.Eventuallyf(t, do, 50*time.Second, 500*time.Millisecond, "controller didn't become ready")
	})

	t.Run("graphql query endpoint is reachable", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/query", nil)
		require.NoError(t, err)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("graphql playground endpoint is reachable", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/internal/playground", nil)
		require.NoError(t, err)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("insertion of single artist works", func(t *testing.T) {
		str := `{"query": "mutation { upsertArtists(input: [{firstName: \"Rainer\", lastName: \"Ingo\"}]) {firstName lastName}}"}`

		body := strings.NewReader(str)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8080/query", body)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		wanted := `{"data":{"upsertArtists":[{"firstName":"Rainer","lastName":"Ingo"}]}}`

		got, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, wanted, string(got))
	})

	// This should always be the last test in this suite.
	t.Run("metrics endpoint is reachable", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/internal/metrics", nil)
		require.NoError(t, err)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
