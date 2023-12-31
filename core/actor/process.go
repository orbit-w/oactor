package actor

// IProcess A process implemented using this module has a
// standard set of interface functions and includes functionality
// for tracing and error reporting
// Let the caller use secure messaging methods to operate the target GenServer
type IProcess interface {
	Cast(pid PID, msg any)
	CastSystem(pid PID, msg any)
	Call(pid PID, msg any) (any, error)
	Stop()
}
