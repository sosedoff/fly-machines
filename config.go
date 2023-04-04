package machines

type Config struct {
	Image       string            `json:"image"`
	Env         map[string]string `json:"env"`
	Init        *InitConfig       `json:"init"`
	Restart     *RestartConfig    `json:"restart"`
	Guest       *GuestConfig      `json:"guest"`
	AutoDestroy bool              `json:"auto_destroy"`
	Schedule    Schedule          `json:"schedule"`
}

type InitConfig struct {
	Exec       string   `json:"exec,omitempty"`
	Cmd        []string `json:"cmd,omitempty"`
	Entrypoint []string `json:"entrypoint,omitempty"`
	TTY        bool     `json:"tty"`
}

type RestartConfig struct {
	Policy     RestartPolicy `json:"policy,omitempty"`
	MaxRetries int           `json:"max_retries,omitempty"`
}

type GuestConfig struct {
	CPUKind CPUKind `json:"cpu_kind"`
	CPUs    uint    `json:"cpus"`
	Memory  uint    `json:"memory_mb"`
}
