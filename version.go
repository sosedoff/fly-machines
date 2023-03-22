package machines

const (
	Version = "0.0.1"
)

func ClientVersion() string {
	return "fly-machines/" + Version
}
