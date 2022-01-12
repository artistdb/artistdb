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

type graphQLResponse struct {
	Data   data           `json:"data"`
	Errors []graphQLError `json:"errors"`
}

type data struct {
	GetArtists         []model.Artist   `json:"getArtists"`
	UpsertArtists      []model.Artist   `json:"upsertArtists"`
	DeleteArtistByID   bool             `json:"deleteArtistByID"`
	UpsertLocations    []model.Location `json:"upsertLocations"`
	DeleteLocationByID bool             `json:"deleteLocationByID"`
}

type graphQLError struct {
	Message string `json:"message"`
}

var httpClient = &http.Client{}

func TestServerIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

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

	t.Run("version endpoint is reachable", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/internal/version", nil)
		require.NoError(t, err)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, []byte("dev"), body)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
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

	t.Run("pprof endpoint is reachable", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/internal/pprof/", nil)
		require.NoError(t, err)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("test artists endpoints", func(t *testing.T) {
		var testID string
		t.Run("insertion of single artist works", func(t *testing.T) {
			str := `{"query": "mutation { upsertArtists(input: [{firstName:\"Bob\",lastName:\"Ross\",artistName:\"BBR\",pronouns:[\"they\",\"them\"],dateOfBirth:1637830936,placeOfBirth:\"Space, Sachsen-Anhalt\",nationality:\"none\",language:\"peace\",facebook:\"meta ;)\",instagram:\"da_real_bob_ross\",bandcamp:\"bandcamp.com/babitorossi\",bioGer:\"bob ross malt so schön!!!!\",bioEn:\"i like so much to draw with bobby\"}]) { id firstName lastName}}"}`

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 0, result.Errors)

			require.Len(t, result.Data.UpsertArtists, 1)
			assert.NotEmpty(t, result.Data.UpsertArtists[0].ID)
			assert.Equal(t, "Bob", result.Data.UpsertArtists[0].FirstName)
			assert.Equal(t, "Ross", result.Data.UpsertArtists[0].LastName)

			testID = result.Data.UpsertArtists[0].ID
		})

		t.Run("retrieval of single artist by ID works", func(t *testing.T) {
			str := fmt.Sprintf(`{"query": "{getArtists(input: [{id: \"%s\"}]){ id, lastName, artistName}}"}`, testID)

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 0, result.Errors)

			require.Len(t, result.Data.GetArtists, 1)
			assert.Equal(t, testID, result.Data.GetArtists[0].ID)
		})

		t.Run("retrieval of single artist by last name works", func(t *testing.T) {
			str := fmt.Sprintf(`{"query": "{getArtists(input: [{lastName: \"%s\"}]){id lastName artistName}}"}`, "Ross")

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 0, result.Errors)

			require.Len(t, result.Data.GetArtists, 1)
			assert.Equal(t, "Ross", result.Data.GetArtists[0].LastName)
		})

		t.Run("retrieval of single artist by artist name works", func(t *testing.T) {
			str := fmt.Sprintf(`{"query": "{getArtists(input: [{lastName: \"%s\"}]){id lastName artistName}}"}`, "Ross")

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 0, result.Errors)

			require.Len(t, result.Data.GetArtists, 1)
			assert.Equal(t, "BBR", *result.Data.GetArtists[0].ArtistName)
		})

		t.Run("Retrieval with invalid ID throws error", func(t *testing.T) {
			str := fmt.Sprintf(`{"query": "{getArtists(input: [{id: \"%s\"}]){ id, lastName, artistName}}"}`, "bogusßß")

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 1, result.Errors)

			require.Len(t, result.Data.GetArtists, 0)
			assert.Contains(t, result.Errors[0].Message, "resource not found")
		})

		t.Run("deletion of single artist works", func(t *testing.T) {
			str := fmt.Sprintf(`{"query": "mutation { deleteArtistByID(id: \"%s\")}"}`, testID)

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 0, result.Errors)

			assert.Equal(t, true, result.Data.DeleteArtistByID)
		})
	})

	t.Run("test locations endpoints", func(t *testing.T) {
		var testID string

		t.Run("insertion of single location works", func(t *testing.T) {
			str := `{"query": 
			"mutation { upsertLocations(input: [{name: \"Tille\"}]) { id name }}"}`

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 0, result.Errors)

			require.Len(t, result.Data.UpsertLocations, 1)
			assert.NotEmpty(t, result.Data.UpsertLocations[0].ID)
			assert.Equal(t, "Tille", result.Data.UpsertLocations[0].Name)

			testID = result.Data.UpsertLocations[0].ID
		})

		t.Run("deletion of single location works", func(t *testing.T) {
			str := fmt.Sprintf(`{"query": "mutation { deleteLocationByID(input: \"%s\")}"}`, testID)

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 0, result.Errors)

			assert.Equal(t, true, result.Data.DeleteLocationByID)
		})
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

func graphQuery(t *testing.T, ctx context.Context, query string) graphQLResponse {
	body := strings.NewReader(query)

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

	var result graphQLResponse
	require.NoError(t, json.Unmarshal(got, &result))

	return result
}
