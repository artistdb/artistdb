package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/obitech/artist-db/internal/database/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type gqlresp struct {
	Data data `json:"data"`
}

type data struct {
	Data          string
	UpsertArtists []model.Artist `json:"upsertArtists"`
}

func TestApiIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpClient := &http.Client{}

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

		got, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		fmt.Println(string(got))

		gqlresp := gqlresp{Data: data{UpsertArtists: []model.Artist{}}}

		marshalErr := json.Unmarshal(got, &gqlresp)
		require.NoError(t, marshalErr)

		fmt.Println(gqlresp)

		assert.Equal(t, "Rainer", gqlresp.Data.UpsertArtists[0].FirstName)
		assert.Equal(t, "Ingo", gqlresp.Data.UpsertArtists[0].LastName)
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
