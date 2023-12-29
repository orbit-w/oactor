package actor

import (
	"github.com/orbit-w/golib/modules/mailbox"
	"sync/atomic"
)

const (
	stateAlive int32 = iota
	stateRestarting
	stateStopping
	stateStopped
)

type OActor struct {
	state    atomic.Int32
	deadFlag atomic.Uint32 // 0 | 1, 1:代表已进入立即注销流程
	behavior IActorBehavior
	mailbox  mailbox.IMailbox //bounded mailbox
}

func (oa *OActor) Cast(_ PID, msg any) {
	oa.mailbox.Push(msg)
}

func (oa *OActor) CastSystem(_ PID, msg any) {
	oa.mailbox.PushSystemMsg(msg)
}

// Stop will tell actor to stop after processing current user messages in mailbox
func (oa *OActor) Stop() {
	oa.mailbox.Push(stopMsg)
}

// Shutdown will stop actor immediately regardless of existing user messages in mailbox.
func (oa *OActor) Shutdown() {
	oa.die()
	oa.mailbox.PushSystemMsg(&SystemStop{})
}

func (oa *OActor) InvokeMsg(message any) {
	switch message.(type) {
	case *Stop:
		oa.Shutdown()
	default:

	}
}

func (oa *OActor) InvokeSysMsg(message any) {
	switch message.(type) {
	case *SystemStop:
		oa.handleStop()
	}
}

func (oa *OActor) handleStop() {
	if oa.state.Load() >= stateStopping {
		return
	}
	oa.state.Store(stateStopping)

	oa.InvokeMsg(stoppingMsg)

	oa.finalizeStop()
}

func (oa *OActor) finalizeStop() {
	oa.InvokeMsg(stoppedMsg)
	oa.state.Store(stateStopped)
}

func (oa *OActor) die() {
	oa.deadFlag.Store(1)
}

func (oa *OActor) dead() bool {
	return oa.deadFlag.Load() == 1
}
