package remote

import (
	"context"
	"github.com/orbit-w/oactor/core/actor"
	"github.com/orbit-w/rpc_transport/rpc"
	"go.uber.org/zap"
)

type Remote struct {
	engine  *actor.Engine
	nodeId  string
	address string
	connMap *ConnMap
	codec   Codec
	logger  *zap.Logger
}

var remote *Remote

func setRemote(r *Remote) {
	if r == nil {
		panic("remote invalid")
	}
	remote = r
}

func NewRemote(e *actor.Engine) (*Remote, error) {
	remote = &Remote{
		engine: e,
		nodeId: e.GetNodeId(),
	}

	remote.connMap = NewConnMap(remote)
	e.RegRemoteHandler(remote.newRemoteProcess)

	if err := rpc.Serve(e.GetNodeId(), func(req rpc.IRequest) error {
		return handleReq(remote, req)
	}); err != nil {
		return nil, err
	}
	return remote, nil
}

func (r *Remote) NodeId() string {
	return r.nodeId
}

func (r *Remote) SendMsg(pid, sender *actor.PID, msg any) error {
	bytes, err := r.codec.Encode(pid, sender, msg)
	if err != nil {
		return err
	}
	return r.connMap.Get(pid).Shoot(bytes)
}

func (r *Remote) Call(ctx context.Context, pid, sender *actor.PID, msg any) (any, error) {
	bytes, err := r.codec.Encode(pid, sender, msg)
	if err != nil {
		return nil, err
	}
	in, err := r.connMap.Get(pid).Call(ctx, bytes)
	if err != nil {
		return nil, err
	}

	return Deserialize(in)
}

func (r *Remote) newRemoteProcess(self *actor.PID) (actor.IProcess, bool) {
	p := newProcess(self, r)
	return p, true
}

func handleReq(r *Remote, in rpc.IRequest) error {
	req := newRequest()

	var err error
	req.Receiver, req.Sender, req.Message, err = r.codec.Decode(in.Data())
	if err != nil {
		return err
	}

	req.Category = in.Category()
	req.Resp = in
	req.Do()
	return nil
}
