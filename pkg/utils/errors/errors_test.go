package pkg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppError_ToHTTPError(t *testing.T) {
	err := NewDomainErrorSimple("CODE", "msg", 400)
	require.Equal(t, "CODE", err.ToHTTPError().Code)
	require.Equal(t, "msg", err.ToHTTPError().Message)
}

func TestToHTTPError_WithAppError(t *testing.T) {
	err := NewApplicationError("C", "m", errors.New("inner"), 500)
	require.Equal(t, "C", ToHTTPError(err).Code)
}

func TestToHTTPError_GenericError(t *testing.T) {
	require.Equal(t, "INTERNAL_ERROR", ToHTTPError(errors.New("x")).Code)
}
