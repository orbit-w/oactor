package actor

type Engine struct {
	address           string
	nodeId            string
	register          *Register
	deadLetterProcess IProcess //死信处理器

	//actor 远程通信处理回调, 可做remote send 重要前置检查
	//默认为 newRemoteProcess 事例化远程虚拟消息通道，不能保证message必达
	remoteHandler []func(*PID) (IProcess, bool)
}

func GEngine() *Engine {
	return gEngine
}

var gEngine *Engine

func NewEngine(conf *Config) *Engine {
	gEngine = new(Engine)
	gEngine.register = NewRegister(128)
	gEngine.remoteHandler = make([]func(*PID) (IProcess, bool), 0)
	return gEngine
}

func (e *Engine) Register() *Register {
	return e.register
}

func (e *Engine) GetNodeId() string {
	return e.nodeId
}

func (e *Engine) LocalAddress() string {
	return e.address
}

func (e *Engine) RegRemoteHandler(h func(*PID) (IProcess, bool)) {
	if h == nil {
		return
	}
	e.remoteHandler = append(e.remoteHandler, h)
}
