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

func (r *Process) Cast(pid *actor.PID, msg any) {
	_ = r.remote.SendMsg(pid, r.self, msg)
}

func (r *Process) CastSystem(_ *actor.PID, msg any) {
	switch msg.(type) {

	}
}

func (r *Process) Call(ctx context.Context, pid *actor.PID, msg any) (any, error) {
	r.remote.Call(ctx, pid, r.self, msg)
	return nil, nil
}

func (r *Process) Stop() {

}
