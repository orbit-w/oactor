package actor

import (
	"sync/atomic"
	"unsafe"
)

// NewPID returns a new instance of the PID struct.
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
		if l, ok := (*p).(*OActor); ok && l.dead() {
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
	return pid.Address != gEngine.localAddress()
}

func FindProcess(pid *PID) (IProcess, bool) {
	if !pid.IsLocal() {
		if ref, ok := GEngine().remoteHandler(pid); ok {
			return ref, true
		}
		return GEngine().deadLetterProcess, false
	}

	ref, exists := gEngine.Register().Get(pid)
	if exists {
		return ref, exists
	}

	return GEngine().deadLetterProcess, false
}
