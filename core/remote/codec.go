package remote

import (
	"github.com/gogo/protobuf/proto"
	"github.com/orbit-w/oactor/core/actor"
)

type Codec struct{}

// Encode 编码
func (c Codec) Encode(pid, sender *actor.PID, msg any) ([]byte, error) {
	var (
		body []byte
		err  error
	)
	if msg != nil {
		body, err = Serialize(msg)
		if err != nil {
			return nil, err
		}
	}

	me := &MessageEnvelope{
		Target: pid,
		Sender: sender,
		Data:   body,
	}
	return proto.Marshal(me)
}

func (c Codec) Decode(in []byte) (target, sender *actor.PID, msg any, err error) {
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
