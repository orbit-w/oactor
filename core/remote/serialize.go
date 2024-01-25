package remote

import (
	"errors"
	"github.com/orbit-w/golib/bases/packet"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var (
	ErrMustProtoMessage = errors.New("must be protobuf v2 message")

	defISerializer = ProtoSerializer{}
)

// ISerializer 传输协议消息序列化解释器
type ISerializer interface {
	Serialize(msg any) ([]byte, error)                   // 序列化
	Deserialize(fullName string, in []byte) (any, error) //反序列化
	FullName(msg any) (string, error)                    //获取传输协议消息体全剧唯一完整名称
}

func Serialize(msg any) ([]byte, error) {
	if msg == nil {
		return []byte{}, nil
	}

	s := defISerializer
	name, err := s.FullName(msg)
	if err != nil {
		return nil, err
	}

	body, err := s.Serialize(msg)
	if err != nil {
		return nil, err
	}

	writer := packet.Writer()
	defer writer.Return()
	writer.WriteString(name)
	writer.Write(body)
	return writer.Copy(), nil
}

func Deserialize(in []byte) (any, error) {
	if len(in) == 0 {
		return nil, nil
	}

	reader := packet.Reader(in)
	defer reader.Return()
	ret, err := reader.ReadBytes()
	if err != nil {
		return nil, err
	}

	name := string(ret)
	s := defISerializer
	return s.Deserialize(name, reader.Remain())
}

type ProtoSerializer struct{}

func (ps ProtoSerializer) Serialize(msg any) ([]byte, error) {
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return nil, ErrMustProtoMessage
	}

	bytes, err := proto.Marshal(pbMsg)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (ps ProtoSerializer) Deserialize(fullName string, bytes []byte) (any, error) {
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(fullName))
	if err != nil {
		return nil, err
	}

	res := msgType.New().Interface()

	err = proto.Unmarshal(bytes, res)
	return res, err
}

func (ps ProtoSerializer) FullName(msg any) (string, error) {
	message, ok := msg.(proto.Message)
	if !ok {
		return "", ErrMustProtoMessage
	}

	return string(proto.MessageName(message)), nil
}
