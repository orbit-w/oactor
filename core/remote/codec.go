package remote

import (
	"github.com/gogo/protobuf/proto"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/oactor/core/actor"
)

type Codec struct {
}

func (c Codec) Encode(pid, sender *actor.PID, msg proto.Message) (packet.IPacket, error) {
	writer := packet.Writer()
	var (
		body []byte
		err  error
	)
	if msg != nil {
		body, err = proto.Marshal(msg)
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
