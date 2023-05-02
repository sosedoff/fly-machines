package machines

type Config struct {
	Image       string                 `json:"image"`
	Env         map[string]string      `json:"env"`
	Init        *InitConfig            `json:"init"`
	Restart     *RestartConfig         `json:"restart"`
	Guest       *GuestConfig           `json:"guest"`
	AutoDestroy bool                   `json:"auto_destroy"`
	Schedule    Schedule               `json:"schedule,omitempty"`
	Services    []ServiceConfig        `json:"services,omitempty"`
	Checks      map[string]CheckConfig `json:"checks,omitempty"`
	Metadata    map[string]string      `json:"metadata,omitempty"`
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

type ServiceConfig struct {
	Protocol     Protocol           `json:"protocol"`                // Networking protocol (TCP/HTTP)
	Concurrency  *ConcurrencyConfig `json:"concurrency,omitempty"`   // Load balancing concurrency settings
	InternalPort uint               `json:"internal_port,omitempty"` // Port the machine VM listens on
	Ports        []PortConfig       `json:"ports,omitempty"`         // Service's ports and associated handler
}

type ConcurrencyConfig struct {
	Type      string `json:"type,omitempty"`
	SoftLimit *int   `json:"soft_limit,omitempty"`
	HardLimit *int   `json:"hard_limit,omitempty"`
}

type PortConfig struct {
	Port       *int     `json:"port,omitempty"`
	Handlers   []string `json:"handlers,omitempty"`
	ForceHTTPS *bool    `json:"force_https,omitempty"`
}

type CheckConfig struct {
	Type          string        `json:"type"`
	Protocol      string        `json:"protocol,omitempty"`
	Port          int           `json:"port,omitempty"`
	Interval      string        `json:"interval"`
	Timeout       string        `json:"timeout"`
	Method        string        `json:"method,omitempty"`
	Path          string        `json:"path,omitempty"`
	TLSSkipVerify *bool         `json:"tls_skip_verify,omitempty"`
	Headers       []CheckHeader `json:"headers,omitempty"`
}

type CheckHeader struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}
