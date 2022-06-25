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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/obitech/artist-db/graph/model"
	"github.com/obitech/artist-db/internal/conversion"
)

type graphQLResponse struct {
	Data   data           `json:"data"`
	Errors []graphQLError `json:"errors"`
}

type data struct {
	GetArtists       []model.Artist `json:"getArtists"`
	UpsertArtists    []model.Artist `json:"upsertArtists"`
	DeleteArtistByID bool           `json:"deleteArtistByID"`

	GetLocations       []model.Location `json:"getLocations"`
	UpsertLocations    []string         `json:"upsertLocations"`
	DeleteLocationByID bool             `json:"deleteLocationByID"`

	UpsertEvents    []string      `json:"upsertEvents"`
	GetEvents       []model.Event `json:"getEvents"`
	DeleteEventByID bool          `json:"deleteEventByID"`
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
			str := `{"query": "mutation { upsertArtists(input: [{firstName:\"Bob\",lastName:\"Ross\",artistName:\"BBR\",pronouns:[\"they\",\"them\"],dateOfBirth:1637830936,placeOfBirth:\"Space, Sachsen-Anhalt\",nationality:\"none\",language:\"peace\",facebook:\"meta ;)\",instagram:\"da_real_bob_ross\",bandcamp:\"bandcamp.com/babitorossi\",bioGer:\"bob ross malt so schön!!!!\",bioEn:\"i like so much to draw with bobby\",email:\"foo@bar.com\"}]) { id firstName lastName email}}"}`

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 0, result.Errors)

			require.Len(t, result.Data.UpsertArtists, 1)
			assert.NotEmpty(t, result.Data.UpsertArtists[0].ID)
			assert.Equal(t, "Bob", result.Data.UpsertArtists[0].FirstName)
			assert.Equal(t, "Ross", result.Data.UpsertArtists[0].LastName)
			assert.Equal(t, "foo@bar.com", conversion.String(result.Data.UpsertArtists[0].Email))

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
			"mutation { upsertLocations(input: [{name: \"Tille\"}])}"}`

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 0, result.Errors)

			require.Len(t, result.Data.UpsertLocations, 1)
			assert.NotEmpty(t, result.Data.UpsertLocations[0])

			_, err := uuid.Parse(result.Data.UpsertLocations[0])
			require.NoError(t, err)

			testID = result.Data.UpsertLocations[0]
		})

		t.Run("retrieval of single location by ID works", func(t *testing.T) {
			str := fmt.Sprintf(`{"query": "{getLocations(input: [{id: \"%s\"}]){id name}}"}`, testID)

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 0, result.Errors)

			require.Len(t, result.Data.GetLocations, 1)
			assert.Equal(t, testID, result.Data.GetLocations[0].ID)
		})

		t.Run("retrieval of single location by name works", func(t *testing.T) {
			str := fmt.Sprintf(`{"query": "{getLocations(input: [{name: \"%s\"}]){id name}}"}`, "Tille")

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 0, result.Errors)

			require.Len(t, result.Data.GetLocations, 1)
			assert.Equal(t, "Tille", result.Data.GetLocations[0].Name)
		})

		t.Run("Retrieval with invalid ID throws error", func(t *testing.T) {
			str := fmt.Sprintf(`{"query": "{getLocations(input: [{id: \"%s\"}]){id name}}"}`, "bogusßß")

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 1, result.Errors)

			require.Len(t, result.Data.GetLocations, 0)
			assert.Contains(t, result.Errors[0].Message, "resource not found")
		})

		t.Run("deletion of single location works", func(t *testing.T) {
			str := fmt.Sprintf(`{"query": "mutation { deleteLocationByID(input: \"%s\")}"}`, testID)

			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 0, result.Errors)

			assert.Equal(t, true, result.Data.DeleteLocationByID)
		})
	})

	t.Run("test events endpoints", func(t *testing.T) {
		t.Run("insertion of single, simple event works", func(t *testing.T) {
			var (
				id        string
				name      = "Ballern"
				startTime = 1637830936
			)

			str := fmt.Sprintf(`{"query": "mutation { upsertEvents(input: [{name: \"%s\", startTime:%d}])}"}`, name, startTime)
			result := graphQuery(t, ctx, str)
			require.Len(t, result.Errors, 0, result.Errors)

			require.Len(t, result.Data.UpsertEvents, 1)
			id = result.Data.UpsertEvents[0]
			assert.NotEmpty(t, id)

			_, err := uuid.Parse(id)
			require.NoError(t, err)

			t.Run("retrieval works", func(t *testing.T) {
				str := fmt.Sprintf(`{"query": "{getEvents(input: [{id: \"%s\"}]){ id, name, startTime}}"}`, id)
				result := graphQuery(t, ctx, str)
				require.Len(t, result.Errors, 0, result.Errors)

				assert.Len(t, result.Data.GetEvents, 1)
				assert.Equal(t, name, result.Data.GetEvents[0].Name)
				assert.Equal(t, startTime, *result.Data.GetEvents[0].StartTime)
			})
		})

		t.Run("insertion of single event+location works", func(t *testing.T) {
			var (
				locID   string
				locName = "Bierkönig"
			)

			t.Run("insertion of single event+location without locationID throws error", func(t *testing.T) {
				str := `{"query": "mutation { upsertEvents(input: [{name: \"Ballern\", startTime:1637830936, locationID: \"foo\"}])}"}`
				result := graphQuery(t, ctx, str)
				require.Len(t, result.Errors, 1, result.Errors)
				assert.Contains(t, result.Errors[0].Message, "invalid UUID")
			})

			t.Run("create new location", func(t *testing.T) {
				str := fmt.Sprintf(`{"query": "mutation { upsertLocations(input: [{name: \"%s\"}])}"}`, locName)

				result := graphQuery(t, ctx, str)
				require.Len(t, result.Errors, 0, result.Errors)

				require.Len(t, result.Data.UpsertLocations, 1)
				assert.NotEmpty(t, result.Data.UpsertLocations[0])

				_, err := uuid.Parse(result.Data.UpsertLocations[0])
				require.NoError(t, err)

				locID = result.Data.UpsertLocations[0]
			})

			t.Run("create event+location", func(t *testing.T) {
				var (
					id        string
					name      = "Ballern 2"
					startTime = 1637830936
				)
				str := fmt.Sprintf(`{"query": "mutation { upsertEvents(input: [{name: \"%s\", startTime:%d, locationID: \"%s\"}])}"}`, name, startTime, locID)
				result := graphQuery(t, ctx, str)

				require.Len(t, result.Errors, 0, result.Errors)

				require.Len(t, result.Data.UpsertEvents, 1)
				id = result.Data.UpsertEvents[0]
				assert.NotEmpty(t, id)

				_, err := uuid.Parse(id)
				require.NoError(t, err)

				t.Run("retrieval works", func(t *testing.T) {
					str := fmt.Sprintf(`{"query": "{getEvents(input: [{id: \"%s\"}]){ id, name, startTime, location {id, name}}}"}`, id)
					result := graphQuery(t, ctx, str)
					require.Len(t, result.Errors, 0, result.Errors)

					assert.Len(t, result.Data.GetEvents, 1)
					assert.Equal(t, name, result.Data.GetEvents[0].Name)
					assert.Equal(t, startTime, *result.Data.GetEvents[0].StartTime)

					assert.Equal(t, locID, result.Data.GetEvents[0].Location.ID)
					assert.Equal(t, locName, result.Data.GetEvents[0].Location.Name)
				})

				t.Run("delete works", func(t *testing.T) {
					str := fmt.Sprintf(`{"query": "mutation { deleteEventByID(input: \"%s\")}"}`, id)

					result := graphQuery(t, ctx, str)
					require.Len(t, result.Errors, 0)

					assert.Equal(t, true, result.Data.DeleteEventByID)

					t.Run("verify event is deleted", func(t *testing.T) {
						str := fmt.Sprintf(`{"query": "{getEvents(input: [{id: \"%s\"}]){ id, name, startTime, location {id, name}}}"}`, id)
						result := graphQuery(t, ctx, str)
						require.Len(t, result.Errors, 1, result.Errors)
						require.Contains(t, result.Errors[0].Message, "resource not found")
					})

					t.Run("location should still exist", func(t *testing.T) {
						str := fmt.Sprintf(`{"query": "{getLocations(input: [{id: \"%s\"}]){id}}"}`, locID)

						result := graphQuery(t, ctx, str)
						require.Len(t, result.Errors, 0, result.Errors)

						require.Len(t, result.Data.GetLocations, 1)
						assert.Equal(t, locID, result.Data.GetLocations[0].ID)
					})
				})
			})
		})

		t.Run("insertion of single event+invited artist works", func(t *testing.T) {
			var (
				artistID        string
				artistFirstName = "DJ"
				artistLastName  = "Robin"
			)

			t.Run("insertion of single event+artist without artistID throws error", func(t *testing.T) {
				str := `{"query": "mutation { upsertEvents(input: [{name: \"Ballern 3\", startTime:1637830936, invitedArtists: [{ id: \"foo\", confirmed: false}]}])}"}`
				result := graphQuery(t, ctx, str)
				require.Len(t, result.Errors, 1, result.Errors)
				assert.Contains(t, result.Errors[0].Message, "invalid UUID")
			})

			t.Run("create new artist", func(t *testing.T) {
				str := fmt.Sprintf(`{"query": "mutation { upsertArtists(input: [{firstName:\"%s\",lastName:\"%s\"}]) {id} }"}`, artistFirstName, artistLastName)

				result := graphQuery(t, ctx, str)
				require.Len(t, result.Errors, 0, result.Errors)

				require.Len(t, result.Data.UpsertArtists, 1)
				assert.NotEmpty(t, result.Data.UpsertArtists[0].ID)

				artistID = result.Data.UpsertArtists[0].ID
			})

			t.Run("create event+artist", func(t *testing.T) {
				var (
					id              string
					eventName       = "Ballern 3"
					eventStart      = 1637830936
					artistConfirmed = false
				)
				str := fmt.Sprintf(`{"query": "mutation { upsertEvents(input: [{name: \"%s\", startTime:%d, invitedArtists: [{ id: \"%s\", confirmed: %t}]}])}"}`, eventName, eventStart, artistID, artistConfirmed)
				result := graphQuery(t, ctx, str)
				require.Len(t, result.Errors, 0, result.Errors)

				require.Len(t, result.Data.UpsertEvents, 1)
				id = result.Data.UpsertEvents[0]
				assert.NotEmpty(t, id)

				_, err := uuid.Parse(id)
				require.NoError(t, err)

				t.Run("retrieval works", func(t *testing.T) {
					str := fmt.Sprintf(`{"query": "{getEvents(input: [{id: \"%s\"}]){ id, name, startTime, artists {artist {id, firstName, lastName}, confirmed}}}"}`, id)
					result := graphQuery(t, ctx, str)
					require.Len(t, result.Errors, 0, result.Errors)

					assert.Len(t, result.Data.GetEvents, 1)
					assert.Equal(t, eventName, result.Data.GetEvents[0].Name)
					assert.Equal(t, eventStart, *result.Data.GetEvents[0].StartTime)

					assert.Len(t, result.Data.GetEvents[0].Artists, 1)
					assert.Equal(t, artistID, result.Data.GetEvents[0].Artists[0].Artist.ID)
					assert.Equal(t, artistFirstName, result.Data.GetEvents[0].Artists[0].Artist.FirstName)
					assert.Equal(t, artistLastName, result.Data.GetEvents[0].Artists[0].Artist.LastName)
					assert.Equal(t, artistConfirmed, result.Data.GetEvents[0].Artists[0].Confirmed)
				})
			})
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
