package remote

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/oactor/core/actor"
)

// ICodec defines the interface that decode/encode payload.
type ICodec interface {
	Encode(msg any) ([]byte, error)
	Decode(buf []byte, r any) error
}

var (
	defaultCodec ICodec
)

func init() {
	//目前只支持 proto 协议
	if defaultCodec == nil {
		defaultCodec = new(ProtoCodec)
	}
}

type ProtoCodec struct{}

func (pi *ProtoCodec) Encode(msg any) ([]byte, error) {
	m, ok := msg.(proto.Message)
	if !ok {
		return nil, errors.New("message not proto.Message")
	}

	body, err := proto.Marshal(m)
	return body, err
}

func (pi *ProtoCodec) Decode(buf []byte, r any) error {
	m, ok := r.(proto.Message)
	if !ok {
		return errors.New("message not proto.Message")
	}
	return proto.Unmarshal(buf, m)
}

type Codec struct{}

func (c Codec) Encode(pid, sender *actor.PID, msg any) (packet.IPacket, error) {
	var (
		body   []byte
		err    error
		writer = packet.Writer()
	)
	if msg != nil {
		body, err = defaultCodec.Encode(msg)
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
