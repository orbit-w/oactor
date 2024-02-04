package actor

type IActorBehavior interface {
	HandleCall(ctx any) (any, error)
}
