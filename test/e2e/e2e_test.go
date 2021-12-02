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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/obitech/artist-db/graph/model"
)

type gqlresp struct {
	Data data `json:"data"`
}

type data struct {
	GetArtists       []model.Artist `json:"getArtists"`
	UpsertArtists    []model.Artist `json:"upsertArtists"`
	DeleteArtistByID bool           `json:"deleteArtistByID"`
}

func TestApiIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	httpClient := &http.Client{}

	var testID string

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
		str := `{"query": 
			"mutation { upsertArtists(input: [{firstName:\"Bob\",lastName:\"Ross\",artistName:\"BBR\",pronouns:[\"they\",\"them\"],dateOfBirth:1637830936,placeOfBirth:\"Space, Sachsen-Anhalt\",nationality:\"none\",language:\"peace\",facebook:\"meta ;)\",instagram:\"da_real_bob_ross\",bandcamp:\"bandcamp.com/babitorossi\",bioGer:\"bob ross malt so schön!!!!\",bioEn:\"i like so much to draw with bobby\"}]) { id firstName lastName}}"}`

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

		var result gqlresp

		unmarshalErr := json.Unmarshal(got, &result)
		require.NoError(t, unmarshalErr)

		require.Len(t, result.Data.UpsertArtists, 1)
		assert.NotEmpty(t, result.Data.UpsertArtists[0].ID)
		assert.Equal(t, "Bob", result.Data.UpsertArtists[0].FirstName)
		assert.Equal(t, "Ross", result.Data.UpsertArtists[0].LastName)

		testID = result.Data.UpsertArtists[0].ID
	})

	t.Run("retrival of single artist by ID works", func(t *testing.T) {
		str := fmt.Sprintf(`{"query": "{getArtists(input: [{id: \"%s\"}]){ id, lastName, artistName}}"}`, testID)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8080/query", strings.NewReader(str))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		got, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result gqlresp

		unmarshalErr := json.Unmarshal(got, &result)
		require.NoError(t, unmarshalErr)

		require.Len(t, result.Data.GetArtists, 1)
		assert.Equal(t, testID, result.Data.GetArtists[0].ID)
	})

	t.Run("retrival of single artist by last name works", func(t *testing.T) {
		str := fmt.Sprintf(`{"query": "{getArtists(input: [{lastName: \"%s\"}]){id lastName artistName}}"}`, "Ross")

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8080/query", strings.NewReader(str))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		got, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result gqlresp

		unmarshalErr := json.Unmarshal(got, &result)
		require.NoError(t, unmarshalErr)

		require.Len(t, result.Data.GetArtists, 1)
		assert.Equal(t, "Ross", result.Data.GetArtists[0].LastName)
	})

	t.Run("retrival of single artist by artist name works", func(t *testing.T) {
		str := fmt.Sprintf(`{"query": "{getArtists(input: [{lastName: \"%s\"}]){id lastName artistName}}"}`, "Ross")

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8080/query", strings.NewReader(str))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		got, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result gqlresp

		unmarshalErr := json.Unmarshal(got, &result)
		require.NoError(t, unmarshalErr)

		require.Len(t, result.Data.GetArtists, 1)
		assert.Equal(t, "BBR", *result.Data.GetArtists[0].ArtistName)
	})

	t.Run("deletion of single artist works", func(t *testing.T) {
		str := fmt.Sprintf(`{"query": "mutation { deleteArtistByID(id: \"%s\")}"}`, testID)

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

		var result gqlresp

		unmarshalErr := json.Unmarshal(got, &result)
		require.NoError(t, unmarshalErr)

		assert.Equal(t, true, result.Data.DeleteArtistByID)
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
