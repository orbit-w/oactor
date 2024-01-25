package remote

import (
	"github.com/orbit-w/oactor/core/actor"
	"sync"
)

const (
	RpcCategoryCast = int8(iota)
	RpcCategoryCall
)

type IRequest interface {
	Response(msg any) error
	Return()
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
	Sender   *actor.PID
	Receiver *actor.PID
	Resp     IResponse
	Message  any
}

func newRequest() *Request {
	v := reqPool.Get()
	req := v.(*Request)
	return req
}

func (req *Request) Do() {
	req.Receiver.SendMessage(req)
}

func (req *Request) Response(msg any) error {
	data, err := Serialize(msg)
	if err != nil {
		return err
	}

	return req.Resp.Response(data)
}

func (req *Request) Return() {
	req.Sender = nil
	req.Receiver = nil
	req.Message = nil
	req.Resp = nil
	req.Category = 0
	reqPool.Put(req)
}
