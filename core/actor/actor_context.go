package actor

import (
	"github.com/orbit-w/golib/modules/mailbox"
	"sync/atomic"
)

const (
	stateAlive uint32 = iota
	stateRestarting
	stateStopping
	stateStopped
)

type Context struct {
	state    atomic.Uint32
	deadFlag atomic.Uint32 // 0 | 1, 1:代表已进入立即注销流程
	self     *PID
	behavior IActorBehavior
	mailbox  mailbox.IMailbox //bounded mailbox
}

// Stop 模式会抛弃掉当前信箱内的所有消息并安全的Terminate
func (ac *Context) Stop(pid *PID) {
	pid.raf().Stop()
}

// GracefulStop 模式, actor在消耗掉收到sign之前的所有消息之后在安全的Terminate
func (ac *Context) GracefulStop(pid *PID) {
	pid.raf().Cast(pid, gracefulStop)
}

func (ac *Context) InvokeMsg(msg any) {
	ac.handleProcessMsg(msg)
}

func (ac *Context) InvokeSysMsg(msg any) {
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

func (ac *Context) handleProcessMsg(msg any) {
	switch message := msg.(type) {
	case *SystemGracefulStop:
		ac.Stop(ac.self)
	case *Request:
		ac.handleRequest(message)
	}
}

func (ac *Context) handleRequest(req *Request) {
	switch req.category {
	case Call:
		result, err := ac.handleCall(req.msg)
		req.Response(result, err)
	}
}

func (ac *Context) handleCall(msg any) (any, error) {
	switch msg.(type) {
	default:
		return ac.behavior.HandleCall(msg)
	}
}

func (ac *Context) handleShutdown() {
	if ac.state.Load() >= stateStopping {
		return
	}
	ac.state.Store(stateStopping)

	ac.InvokeMsg(stoppingMsg)

	ac.InvokeMsg(stoppedMsg)
	ac.state.Store(stateStopped)
}
