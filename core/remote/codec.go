package remote

import (
	"github.com/gogo/protobuf/proto"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/oactor/core/actor"
)

type Codec struct{}

// EncodeReq 编码
func (c Codec) EncodeReq(pid, sender *actor.PID, msg any) (packet.IPacket, error) {
	var (
		body   []byte
		writer = packet.Writer()
		err    error
	)

	if msg != nil {
		body, err = Serialize(msg)
	}

	me := MessageEnvelope{
		Target: pid,
		Sender: sender,
		Data:   body,
	}
	pack, err := proto.Marshal(&me)
	if err != nil {
		return nil, err
	}
	writer.Write(pack)
	return writer, err
}

func (c Codec) EncodeResp(msg any) (packet.IPacket, error) {
	if msg == nil {
		return nil, nil
	}
	body, err := Serialize(msg)
	return packet.Reader(body), err
}

func (c Codec) DecodeReq(in []byte) (target, sender *actor.PID, msg any, err error) {
	me := new(MessageEnvelope)
	if err = proto.Unmarshal(in, me); err != nil {
		return nil, nil, nil, err
	}
	target = me.Target
	sender = me.Sender
	if me.Data == nil {
		return nil, nil, nil, nil
	}

	if len(me.Data) > 0 {
		msg, err = Deserialize(me.Data)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	return
}

// DecodeResp 解码
func (c Codec) DecodeResp(in []byte) (any, error) {
	return Deserialize(in)
}
