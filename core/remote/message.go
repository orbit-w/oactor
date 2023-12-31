package remote

import (
	"github.com/gogo/protobuf/proto"
	"github.com/orbit-w/oactor/core/actor"
)

const (
	RpcCategoryCast = int8(iota)
	RpcCategoryCall
)

type RpcRequest struct {
	Id       uint64
	Category int8
	Sender   *actor.PID
	Receiver *actor.PID
	Message  proto.Message
}
