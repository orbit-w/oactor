package actor

import (
	"sync/atomic"
	"unsafe"
)

// NewPID returns a new instance of the PID struct.
// Pid ( Process Id ) 消息投递的入口, 一个可序列化的标识符，作为消息投递的入口，将消息发送到参与者的邮箱；
// Pid 可以低消耗的序列化供远程传输
func NewPID(address, id string) *PID {
	return &PID{
		Address: address,
		Id:      id,
	}
}

func (pid *PID) SendMessage(message interface{}) {
	if !pid.IsLocal() {
		//remote
		return
	}
	pid.raf().Cast(pid, message)
}

// TODO：如果actor 发生物理结点迁移，需要重新获取PID
func (pid *PID) raf() IProcess {
	p := (*IProcess)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&pid.p))))
	if p != nil {
		if process, ok := (*p).(*ActorProcess); ok && process.dead() {
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&pid.p)), nil)
		} else {
			return *p
		}
	}

	ref, exists := FindProcess(pid)
	if exists {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&pid.p)), unsafe.Pointer(&ref))
	}
	return ref
}

func (pid *PID) Equal(other *PID) bool {
	if pid != nil && other == nil {
		return false
	}

	return pid.Id == other.Id && pid.Address == other.Address && pid.RequestId == other.RequestId
}

func (pid *PID) IsLocal() bool {
	return pid.Address != gEngine.LocalAddress()
}

func FindProcess(pid *PID) (IProcess, bool) {
	if !pid.IsLocal() {
		e := GEngine()
		//actor 远程通信前置条件检查
		//例如：远端 actor 活跃检查
		for i := range e.remoteHandler {
			h := e.remoteHandler[i]
			if ref, ok := h(pid); ok {
				return ref, ok
			}
		}
		return e.deadLetterProcess, false
	}

	//本地消息通信
	ref, exists := gEngine.Register().Get(pid)
	if exists {
		return ref, exists
	}

	return GEngine().deadLetterProcess, false
}
