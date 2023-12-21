package remote

import (
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/oactor/core/actor"
	"github.com/orbit-w/oactor/core/transport"
	"sync"
)

type Remote struct {
	e       *actor.Engine
	nodeId  string
	pid     *actor.PID
	connMap *ConnMap
}

func (r *Remote) Id() string {
	return r.pid.Id
}

func (r *Remote) SendMsg(pid *actor.PID, msg any) error {
	pack := r.encode(pid, msg)
	defer pack.Return()
	return r.connMap.Get(pid).Write(pack)
}

func (r *Remote) encode(pid *actor.PID, msg any) packet.IPacket {
	return nil
}

type ConnMap struct {
	rw      sync.RWMutex
	remote  *Remote
	connMap map[string]transport.IConn
}

func (rc *ConnMap) Get(t *actor.PID) transport.IConn {
	rc.rw.RLock()
	if conn, ok := rc.connMap[t.Id]; ok {
		rc.rw.RUnlock()
		return conn
	}
	rc.rw.RUnlock()

	return rc.Load(t)
}

func (rc *ConnMap) Load(t *actor.PID) transport.IConn {
	rc.rw.Lock()
	if conn, ok := rc.connMap[t.Id]; ok {
		rc.rw.Unlock()
		return conn
	}

	conn := transport.DialWithOps(t.Address, &transport.DialOption{
		RemoteNodeId:  t.Id,
		CurrentNodeId: rc.remote.Id(),
		DisconnectHandler: func(nodeId string) {
			rc.rw.Lock()
			delete(rc.connMap, nodeId)
			rc.rw.Unlock()
		},
	})
	rc.connMap[t.Id] = conn
	rc.rw.Unlock()
	return conn
}
