package actor

// System messages
type (
	SystemStop         struct{}
	SystemGracefulStop struct{}
	SystemStarted      struct{}
	SystemRestart      struct{}
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

	//系统级消息
	stopMsg      = &SystemStop{}
	gracefulStop = &SystemGracefulStop{}
	startedMsg   = &SystemStarted{}
)
