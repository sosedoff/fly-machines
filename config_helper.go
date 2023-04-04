package machines

// HTTPPortConfig returns default configuration for the HTTP service
func HTTPPortConfig() PortConfig {
	return PortConfig{
		Port:       IntVal(80),
		ForceHTTPS: BoolVal(true),
		Handlers:   []string{"http"},
	}
}

// HTTPSPortConfig returns default configuration for the HTTPS service
func HTTPSPortConfig() PortConfig {
	return PortConfig{
		Port:     IntVal(443),
		Handlers: []string{"tls", "http"},
	}
}

// DefaultWebPortConfigs returns default configuration for the web service running HTTP/HTTPS
func DefaultWebPortConfigs() []PortConfig {
	return []PortConfig{
		HTTPPortConfig(),
		HTTPSPortConfig(),
	}
}
