package actor

import "context"

// IProcess 每个 Actor 都有一个与之关联的 IProcess 实例, IProcess 负责将消息投递给 Actor 以及 Actor 生命周期管理
// IProcess 提供了一个透明的机制来发送消息，使得Local和Remote通信对于 Actor 来说是透明的
// 业务可以自定义 IProcess 来实现不同的业务场景
type IProcess interface {
	Cast(pid *PID, msg any)
	CastSystem(pid *PID, msg any)
	Call(ctx context.Context, pid *PID, msg any) (any, error)
	Stop()
}
