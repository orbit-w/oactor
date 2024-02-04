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

type ActorContext struct {
	self     *PID
	state    atomic.Int32
	deadFlag atomic.Uint32 // 0 | 1, 1:代表已进入立即注销流程
	behavior IActorBehavior
	mailbox  mailbox.IMailbox //bounded mailbox
}

// Stop 模式会抛弃掉当前信箱内的所有消息并安全的Terminate
func (ac *ActorContext) Stop(pid *PID) {
	pid.raf().Stop()
}

// GracefulStop 模式, actor在消耗掉收到sign之前的所有消息之后在安全的Terminate
func (ac *ActorContext) GracefulStop(pid *PID) {
	pid.raf().Cast(pid, gracefulStop)
}

func (ac *ActorContext) InvokeMsg(msg any) {
	ac.handleProcessMsg(msg)
}

func (ac *ActorContext) InvokeSysMsg(msg any) {
	switch msg.(type) {
	case *SystemStop:
		//立即停止，不会消耗掉后续的消息
		ac.handleShutdown()
	case *SystemStarted:
		//激活驱动消息，是Actor在生成或重新启动后收到的第一条消息。
		//如果需要为参与者设置初始状态（例如从数据库加载数据），需要在IActorBehavior中处理SystemStarted消息
		ac.InvokeMsg(msg)
	}
}

func (ac *ActorContext) handleProcessMsg(msg any) {
	switch msg.(type) {
	case *SystemGracefulStop:
		ac.Stop(ac.self)
	default:
		//业务层自定义 IActorBehavior 调用
	}
}

func (ac *ActorContext) handleShutdown() {
	if ac.state.Load() >= stateStopping {
		return
	}
	ac.state.Store(stateStopping)

	ac.InvokeMsg(stoppingMsg)

	ac.InvokeMsg(stoppedMsg)
	ac.state.Store(stateStopped)
}
