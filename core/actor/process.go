package actor

type IProcess interface {
	Cast(msg any)
	CastSystem(msg any)
	Stop()
}
