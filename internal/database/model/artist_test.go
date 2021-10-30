package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewArtist(t *testing.T) {
	require.NotEmpty(t, NewArtist().ID)
}
