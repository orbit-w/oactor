package remote

import (
	"context"
	"github.com/orbit-w/oactor/core/actor"
)

type Process struct {
	self   *actor.PID
	remote *Remote
}

func newProcess(pid *actor.PID, _remote *Remote) *Process {
	return &Process{
		self:   pid,
		remote: _remote,
	}
}

func (p *Process) Cast(pid *actor.PID, msg any) {
	_ = p.remote.SendMsg(pid, p.self, msg)
}

func (p *Process) CastSystem(_ *actor.PID, msg any) {
	switch msg.(type) {

	}
}

func (p *Process) Call(ctx context.Context, pid *actor.PID, msg any) (any, error) {
	return p.remote.Call(ctx, pid, p.self, msg)
}

func (p *Process) Stop() {

}

func (p *Process) GracefulStop() {

}
