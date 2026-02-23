package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	bookmark, err := NewBookmark("test", "value")

	require.NoError(t, err)
	require.Equal(t, "test", bookmark.Title)
	require.Equal(t, "value", bookmark.Value)
}

func TestNewPrepare_Success(t *testing.T) {
	bookmark, err := NewBookmark("test ", " value ")

	require.NoError(t, err)
	require.Equal(t, "test", bookmark.Title)
	require.Equal(t, "value", bookmark.Value)
}

func TestNew_Error(t *testing.T) {
	_, err := NewBookmark("test", "")

	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidValue)

	_, err = NewBookmark("", "")

	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidTitle)
}
