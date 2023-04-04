package machines

type Size string

const (
	SizeSharedCPU1x Size = "shared-cpu-1x"
	SizeSharedCPU2x Size = "shared-cpu_2x"
	SizeSharedCPU4x Size = "shared-cpu-4x"
	SizeSharedCPU8x Size = "shared-cpu-8x"

	SizePerformance1x  Size = "performance-1x"
	SizePerformance2x  Size = "performance-2x"
	SizePerformance4x  Size = "performance-4x"
	SizePerformance8x  Size = "performance-8x"
	SizePerformance16x Size = "performance-16x"
)

type CPUKind string

const (
	CPUKindShared      = "shared"
	CPUKindPerformance = "performance"
)

type State string

const (
	StateCreated    State = "created"
	StateStarting   State = "starting"
	StateStarted    State = "started"
	StateStopping   State = "stopping"
	StateStopped    State = "stopped"
	StateReplacing  State = "replacing"
	StateDestroying State = "destroying"
	StateDestroyed  State = "destroyed"
)

type RestartPolicy string

const (
	RestartPolicyNo        RestartPolicy = "no"
	RestartPolicyOnFailure RestartPolicy = "on-failure"
	RestartPolicyAlways    RestartPolicy = "always"
)

var (
	PolicyConfigRestartNever = &RestartConfig{Policy: RestartPolicyNo, MaxRetries: 0}
	PolicyConfigRestartOnce  = &RestartConfig{Policy: RestartPolicyOnFailure, MaxRetries: 1}
)

type Schedule string

const (
	ScheduleHourly  Schedule = "hourly"
	ScheduleDaily   Schedule = "daily"
	ScheduleWeekly  Schedule = "weekly"
	ScheduleMonthly Schedule = "monthly"
)

type Protocol string

const (
	ProtocolTCP Protocol = "tcp"
	ProtocolUDP Protocol = "udp"
)
