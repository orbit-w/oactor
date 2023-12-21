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
	if pid.Address != engine.localAddress() {
		//remote
		return
	}
	pid.raf().Cast(message)
}

func (pid *PID) raf() IProcess {
	p := (*IProcess)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&pid.p))))
	if p != nil {
		if l, ok := (*p).(*OActor); ok && l.dead() {
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&pid.p)), nil)
		} else {
			return *puukj
		}
	}

	ref, exists := engine.Register().Get(pid)
	if exists {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&pid.p)), unsafe.Pointer(&ref))
		return ref
	}

	return ref
}

func (pid *PID) Equal(other *PID) bool {
	if pid != nil && other == nil {
		return false
	}

	return pid.Id == other.Id && pid.Address == other.Address && pid.RequestId == other.RequestId
}
