package remote

/*
	@Author: orbit-w
	@File: meg_interpreter
	@2024 1月 周日 21:18
*/

import (
	"errors"
	"github.com/gogo/protobuf/proto"
)

var (
	defaultInterpreter IMessageInterpreter
)

// IMessageInterpreter 消息协议解释器
type IMessageInterpreter interface {
	Marshal(msg any) ([]byte, error)
	Unmarshal(buf []byte, r any) error
}

func init() {
	//目前只支持 proto 协议
	if defaultInterpreter == nil {
		defaultInterpreter = new(ProtoInterpreter)
	}
}

func Interpreter() IMessageInterpreter {
	return defaultInterpreter
}

func SetInterpreter(i IMessageInterpreter) {
	if i != nil {
		defaultInterpreter = i
	}
}

type ProtoInterpreter struct{}

func (pi *ProtoInterpreter) Marshal(msg any) ([]byte, error) {
	m, ok := msg.(proto.Message)
	if !ok {
		return nil, errors.New("message not proto.Message")
	}

	body, err := proto.Marshal(m)
	return body, err
}

func (pi *ProtoInterpreter) Unmarshal(buf []byte, r any) error {
	m, ok := r.(proto.Message)
	if !ok {
		return errors.New("message not proto.Message")
	}
	return proto.Unmarshal(buf, m)
}
