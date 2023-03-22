package machines

import (
	"encoding/json"
	"io"
	"net/http"
)

// APIError is a generic error response container
type APIError struct {
	StatusCode   int               `json:"-"`
	ErrorMessage string            `json:"error"`
	Headers      map[string]string `json:"-"`

	rawBody []byte
}

func (err APIError) RawBody() string {
	return string(err.rawBody)
}

func (err APIError) Error() string {
	return err.ErrorMessage
}

func apiErrorFromResponse(resp *http.Response) error {
	apiErr := APIError{
		StatusCode: resp.StatusCode,
		Headers:    map[string]string{},
	}

	// We copy all the headers instead of just FLY_* ones for debugging
	for k, v := range resp.Header {
		apiErr.Headers[k] = v[0]
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	apiErr.rawBody = body

	// Check if body looks like JSON first
	if body[0] == '{' {
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return err
		}
	}

	return apiErr
}
