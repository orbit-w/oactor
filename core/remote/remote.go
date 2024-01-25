package remote

import (
	"context"
	mmrpc "github.com/orbit-w/mmrpc/rpc"
	"github.com/orbit-w/oactor/core/actor"
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

func NewRemote(e *actor.Engine) *Remote {
	remote = &Remote{
		engine: e,
		nodeId: e.GetNodeId(),
	}

	remote.connMap = NewConnMap(remote)
	if err := mmrpc.Serve(e.GetNodeId(), func(req mmrpc.IRequest) error {
		return handleReq(remote, req)
	}); err != nil {

	}
	return remote
}

func (r *Remote) NodeId() string {
	return r.nodeId
}

func (r *Remote) SendMsg(pid, sender *actor.PID, msg any) error {
	req, err := r.codec.EncodeReq(pid, sender, msg)
	if err != nil {
		req.Return()
		return err
	}
	defer req.Return()
	return r.connMap.Get(pid).Shoot(req.Data())
}

func (r *Remote) Call(ctx context.Context, pid, sender *actor.PID, msg any) (any, error) {
	req, err := r.codec.EncodeReq(pid, sender, msg)
	if err != nil {
		req.Return()
		return nil, err
	}
	defer req.Return()
	in, err := r.connMap.Get(pid).Call(ctx, req.Data())
	if err != nil {
		return nil, err
	}

	return r.codec.DecodeResp(in)
}

func handleReq(r *Remote, in mmrpc.IRequest) error {
	req, err := newRequest(r, in)

	if err != nil {
		r.logger.Error("decode req failed", zap.Error(err))
		return err
	}

	doRequest(req)
	return nil
}
