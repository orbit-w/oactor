package actor

// System messages
type (
	SystemStop    struct{}
	SystemStarted struct{}
	SystemRestart struct{}
)

// User auto receive messages
type (
	Stopping struct{}
	Stopped  struct{}
	Stop     struct{}
)

var (
	stoppingMsg = &Stopping{}
	stoppedMsg  = &Stopped{}
	stopMsg     = &Stop{}
)
