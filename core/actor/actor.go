package actor

import (
	"context"
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

func (oa *OActor) Cast(_ *PID, msg any) {
	oa.mailbox.Push(msg)
}

func (oa *OActor) Call(ctx context.Context, pid *PID, msg any) (any, error) {
	return nil, nil
}

func (oa *OActor) CastSystem(_ *PID, msg any) {
	oa.mailbox.PushSystemMsg(msg)
}

// SystemStopMsg will tell actor to stop after processing current user messages in mailbox
func (oa *OActor) Stop() {
	oa.mailbox.PushSystemMsg(stopMsg)
}

// Shutdown will stop actor immediately regardless of existing user messages in mailbox.
func (oa *OActor) Shutdown() {
	oa.die()
	//oa.mailbox.PushSystemMsg(stopMsg)
}

func (oa *OActor) InvokeMsg(message any) {

}

func (oa *OActor) InvokeSysMsg(message any) {
	switch message.(type) {
	case *SystemStop:
		//立即停止，不会消耗掉后续的消息
		oa.handleStop()
	case *SystemStarted:
		//激活驱动消息，是Actor在生成或重新启动后收到的第一条消息。
		//如果需要为参与者设置初始状态（例如从数据库加载数据），需要在IActorBehavior中处理SystemStarted消息
		oa.InvokeMsg(message)
	}
}

func (oa *OActor) handleStop() {
	if oa.state.Load() >= stateStopping {
		return
	}
	oa.state.Store(stateStopping)

	oa.InvokeMsg(stoppingMsg)

	oa.InvokeMsg(stoppedMsg)
	oa.state.Store(stateStopped)
}

func (oa *OActor) die() {
	oa.deadFlag.Store(1)
}

func (oa *OActor) dead() bool {
	return oa.deadFlag.Load() == 1
}
