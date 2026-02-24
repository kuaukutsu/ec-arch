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
	tests := []struct {
		title string
		value string
		err   error
	}{
		{
			title: "test",
			value: "",
			err:   ErrInvalidValue,
		},
		{
			title: "",
			value: "value",
			err:   ErrInvalidTitle,
		},
		{
			title: "",
			value: "",
			err:   ErrInvalidTitle,
		},
	}

	for _, tt := range tests {
		_, err := NewBookmark(tt.title, tt.value)

		require.Error(t, err)
		require.ErrorIs(t, err, tt.err)
	}
}
