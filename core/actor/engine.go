package actor

type Engine struct {
	address  string
	nodeId   string
	register *Register
}

func GEngine() *Engine {
	return engine
}

var engine *Engine

func NewEngine(conf *Config) *Engine {
	engine = new(Engine)
	engine.register = NewRegister(128)
	return engine
}

func (e *Engine) Register() *Register {
	return e.register
}

func (e *Engine) localAddress() string {
	return e.address
}
