package remote

import (
	"github.com/gogo/protobuf/proto"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/oactor/core/actor"
)

type Codec struct {
}

func (c Codec) Encode(pid *actor.PID, msg proto.Message) (packet.IPacket, error) {
	head, err := proto.Marshal(pid)
	if err != nil {
		return nil, err
	}
	writer := packet.Writer()
	writer.WriteBytes(head)

	body, err := proto.Marshal(msg)
	if err == nil {
		writer.WriteBytes(body)
	}
	return writer, err
}
