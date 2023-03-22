package machines

import "fmt"

// type APIResponse struct {
// }

type APIError struct {
	StatusCode   int    `json:"-"`
	ErrorMessage string `json:"error"`
}

func (err APIError) Error() string {
	return fmt.Sprintf("error_code=%d error_message=%s", err.StatusCode, err.ErrorMessage)
}

// type ErrBadRequest struct {
// 	APIError
// }
