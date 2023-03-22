package machines

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_apiErrorFromResponse(t *testing.T) {
	readFixture := func(path string) io.ReadCloser {
		f, err := os.Open("testdata/" + path)
		if err != nil {
			panic(err)
		}
		return io.NopCloser(f)
	}

	t.Run("json error", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 400,
			Header: map[string][]string{
				"fly-request-id": {"req-id"},
				"fly-trace-id":   {"trace-id"},
				"via":            {"server: Fly/620fe63b (2023-03-17)"},
			},
			Body: readFixture("not_found.json"),
		}

		err := apiErrorFromResponse(resp)
		assert.IsType(t, APIError{}, err)

		apiErr := err.(APIError)
		assert.Equal(t, 400, apiErr.StatusCode)
		assert.Equal(t, "invalid machine ID, '12345'", apiErr.ErrorMessage)
		assert.Contains(t, string(apiErr.rawBody), `{"error": "invalid machine ID, '12345'"}`)
		assert.Equal(t, "req-id", apiErr.Headers["fly-request-id"])
		assert.Equal(t, "trace-id", apiErr.Headers["fly-trace-id"])
		assert.Equal(t, "server: Fly/620fe63b (2023-03-17)", apiErr.Headers["via"])
	})

	t.Run("text error", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 400,
			Body:       readFixture("not_found.txt"),
		}

		err := apiErrorFromResponse(resp)
		assert.IsType(t, APIError{}, err)

		apiErr := err.(APIError)
		assert.Equal(t, 400, apiErr.StatusCode)
		assert.Equal(t, "", apiErr.ErrorMessage) // TODO: set error from the body?
		assert.Contains(t, string(apiErr.rawBody), `not found`)
	})
}
