package remote

import (
	"github.com/orbit-w/oactor/core/actor"
	"github.com/orbit-w/rpc_transport/rpc"
	"log"
	"sync"
)

type ConnMap struct {
	rw      sync.RWMutex
	remote  *Remote
	connMap map[string]rpc.IClient
}

func NewConnMap(_remote *Remote) *ConnMap {
	return &ConnMap{
		rw:      sync.RWMutex{},
		remote:  _remote,
		connMap: make(map[string]rpc.IClient),
	}
}

func (rc *ConnMap) Get(t *actor.PID) rpc.IClient {
	rc.rw.RLock()
	if conn, ok := rc.connMap[t.Id]; ok {
		rc.rw.RUnlock()
		return conn
	}
	rc.rw.RUnlock()

	return rc.Load(t)
}

func (rc *ConnMap) Load(t *actor.PID) rpc.IClient {
	rc.rw.Lock()
	defer func() {
		rc.rw.Unlock()
	}()
	if conn, ok := rc.connMap[t.Id]; ok {
		return conn
	}

	conn, err := rpc.Dial(rc.remote.NodeId(), t.NodeId, t.Address, &rpc.DialOption{
		DisconnectHandler: func(nodeId string) {
			rc.rw.Lock()
			delete(rc.connMap, nodeId)
			rc.rw.Unlock()
		},
	})
	if err != nil {
		log.Fatalln("rpc dial failed: ", err.Error())
		return nil
	}

	rc.connMap[t.Id] = conn
	return conn
}
