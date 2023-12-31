package actor

type Engine struct {
	address           string
	nodeId            string
	register          *Register
	deadLetterProcess IProcess
	remoteHandler     func(*PID) (IProcess, bool)
}

func GEngine() *Engine {
	return gEngine
}

var gEngine *Engine

func NewEngine(conf *Config) *Engine {
	gEngine = new(Engine)
	gEngine.register = NewRegister(128)
	return gEngine
}

func (e *Engine) Register() *Register {
	return e.register
}

func (e *Engine) GetNodeId() string {
	return e.nodeId
}

func (e *Engine) localAddress() string {
	return e.address
}
