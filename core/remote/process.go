package remote

import "github.com/orbit-w/oactor/core/actor"

type Process struct {
	self actor.PID
}

func (r *Process) Cast(pid actor.PID, msg any) {
	request := RpcRequest{
		message: msg,
	}
}

func (r *Process) CastSystem(pid actor.PID, msg any) {

}

func (r *Process) Stop() {

}

func newRequest() *RpcRequest {
	return &RpcRequest{}
}
