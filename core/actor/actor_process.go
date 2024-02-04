package actor

import (
	"context"
	"github.com/orbit-w/golib/modules/mailbox"
	"sync/atomic"
)

type ActorProcess struct {
	mailbox  mailbox.IMailbox //bounded mailbox
	deadFlag atomic.Int32
}

func (ap *ActorProcess) Cast(_ *PID, msg any) {
	ap.mailbox.Push(msg)
}

func (ap *ActorProcess) Call(ctx context.Context, pid *PID, msg any) (any, error) {
	return nil, nil
}

func (ap *ActorProcess) CastSystem(_ *PID, msg any) {
	ap.mailbox.PushSystemMsg(msg)
}

// Stop 模式会抛弃掉当前信箱内的所有消息并安全的Terminate
func (ap *ActorProcess) Stop() {
	ap.deadFlag.Store(1)
	ap.mailbox.PushSystemMsg(stopMsg)
}

// GracefulStop 模式, actor在消耗掉收到sign之前的所有消息之后在安全的Terminate
func (ap *ActorProcess) GracefulStop() {
	ap.mailbox.Push(gracefulStop)
}

func (ap *ActorProcess) dead() bool {
	return ap.deadFlag.Load() == 1
}
