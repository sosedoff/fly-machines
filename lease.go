package machines

type Lease struct {
	Nonce     string `json:"nonce"`
	ExpiresAt int64  `json:"expires_at"`
	Owner     string `json:"owner"`
}
