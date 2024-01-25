package remote

import (
	mmrpc "github.com/orbit-w/mmrpc/rpc"
	"github.com/orbit-w/oactor/core/actor"
	"go.uber.org/zap"
	"sync"
)

const (
	RpcCategoryCast = int8(iota)
	RpcCategoryCall
)

type IRequest interface {
	Response(msg any) error
}

type IResponse interface {
	Response(out []byte) error
	Return()
	Category() int8
}

var (
	reqPool = sync.Pool{New: func() any {
		return new(Request)
	}}
)

type Request struct {
	Category int8
	remote   *Remote
	Sender   *actor.PID
	Receiver *actor.PID
	Resp     IResponse
	Message  any
}

func newRequest(r *Remote, in mmrpc.IRequest) (*Request, error) {
	v := reqPool.Get()
	req := v.(*Request)
	target, sender, msg, err := r.codec.DecodeReq(in.Data())
	if err != nil {
		r.logger.Error("decode req failed", zap.Error(err))
		return nil, err
	}

	req.Category = in.Category()
	req.Sender = sender
	req.Receiver = target
	req.Message = msg
	req.Resp = in
	return req, nil
}

func doRequest(req *Request) {
	req.Receiver.SendMessage(req)
}

func (req *Request) Response(msg any) error {
	pack, err := req.remote.codec.EncodeResp(msg)
	defer func() {
		if pack != nil {
			pack.Return()
		}
	}()
	if err != nil {
		return err
	}

	return req.Resp.Response(pack.Data())
}
