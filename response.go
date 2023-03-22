package machines

// type APIResponse struct {
// }

type APIError struct {
	StatusCode   int    `json:"-"`
	ErrorMessage string `json:"error"`
}

func (err APIError) Error() string {
	return err.ErrorMessage
}

// type ErrBadRequest struct {
// 	APIError
// }
