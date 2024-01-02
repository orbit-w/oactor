package remote

import (
	"github.com/gogo/protobuf/proto"
	"github.com/orbit-w/oactor/core/actor"
	"github.com/orbit-w/oactor/core/transport"
	"sync"
)

type Remote struct {
	engine  *actor.Engine
	nodeId  string
	connMap *ConnMap
	codec   Codec
}

var remote *Remote

func NewRemote(e *actor.Engine) *Remote {
	remote = &Remote{
		engine: e,
		nodeId: e.GetNodeId(),
	}

	remote.connMap = NewConnMap(remote)
	return remote
}

func (r *Remote) NodeId() string {
	return r.nodeId
}

func (r *Remote) SendMsg(pid, sender *actor.PID, msg proto.Message) error {
	pack, err := r.codec.Encode(pid, sender, msg)
	if err != nil {
		pack.Return()
		return err
	}
	defer pack.Return()
	return r.connMap.Get(pid).Write(pack)
}

type ConnMap struct {
	rw      sync.RWMutex
	remote  *Remote
	connMap map[string]transport.IConn
}

func NewConnMap(_remote *Remote) *ConnMap {
	return &ConnMap{
		rw:      sync.RWMutex{},
		remote:  _remote,
		connMap: make(map[string]transport.IConn),
	}
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
		CurrentNodeId: rc.remote.NodeId(),
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
