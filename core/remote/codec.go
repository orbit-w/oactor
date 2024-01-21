package remote

import (
	"github.com/gogo/protobuf/proto"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/oactor/core/actor"
)

type Codec struct{}

func (c Codec) Encode(pid, sender *actor.PID, msg any) (packet.IPacket, error) {
	var (
		body   []byte
		err    error
		writer = packet.Writer()
	)
	if msg != nil {
		body, err = Interpreter().Marshal(msg)
		if err != nil {
			return nil, err
		}
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

func (c Codec) Decode(me *MessageEnvelope, data []byte, reader packet.IPacket) error {
	return nil
}
