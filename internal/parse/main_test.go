package parse_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ffigari/stored-strings/internal/parse"
)

func TestParsingWorksWhenNoBodyParamsArePresent(t *testing.T) {
	r, err := http.NewRequest("POST", "http://example.com", nil)
	require.NoError(t, err)

	parsedParams := parse.BodyParams(r)
	require.NotNil(t, parsedParams)
	require.Equal(t, len(parsedParams), 0)
}

func TestParsingWorksWhenOnlyOneParameterIsPresent(t *testing.T) {
	r, err := http.NewRequest(
		"POST",
		"http://example.com",
		bytes.NewBuffer([]byte("username=admin")),
	)
	require.NoError(t, err)

	parsedParams := parse.BodyParams(r)
	require.NotNil(t, parsedParams)

	v, ok := parsedParams["username"]
	assert.True(t, ok)
	assert.Equal(t, v, "admin")
}

func TestParsingWorksWhenMoreThanOneParameterIsPresent(t *testing.T) {
	r, err := http.NewRequest(
		"POST",
		"http://example.com",
		bytes.NewBuffer([]byte("username=admin&password=1234abcd")),
	)
	require.NoError(t, err)

	parsedParams := parse.BodyParams(r)
	require.NotNil(t, parsedParams)

	v0, ok0 := parsedParams["username"]
	assert.True(t, ok0)
	assert.Equal(t, v0, "admin")

	v1, ok1 := parsedParams["password"]
	assert.True(t, ok1)
	assert.Equal(t, v1, "1234abcd")
}
