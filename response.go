package machines

// APIError is a generic error response container
type APIError struct {
	StatusCode   int               `json:"-"`
	ErrorMessage string            `json:"error"`
	Headers      map[string]string `json:"-"`
}

func (err APIError) Error() string {
	return err.ErrorMessage
}
