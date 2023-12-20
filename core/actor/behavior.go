package actor

type IActorBehavior interface {
	HandleMsg(ctx any)
}
