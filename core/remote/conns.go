package remote

import (
	mmrpc "github.com/orbit-w/mmrpc/rpc"
	"github.com/orbit-w/oactor/core/actor"
	"log"
	"sync"
)

type ConnMap struct {
	rw      sync.RWMutex
	remote  *Remote
	connMap map[string]mmrpc.IClient
}

func NewConnMap(_remote *Remote) *ConnMap {
	return &ConnMap{
		rw:      sync.RWMutex{},
		remote:  _remote,
		connMap: make(map[string]mmrpc.IClient),
	}
}

func (rc *ConnMap) Get(t *actor.PID) mmrpc.IClient {
	rc.rw.RLock()
	if conn, ok := rc.connMap[t.Id]; ok {
		rc.rw.RUnlock()
		return conn
	}
	rc.rw.RUnlock()

	return rc.Load(t)
}

func (rc *ConnMap) Load(t *actor.PID) mmrpc.IClient {
	rc.rw.Lock()
	defer func() {
		rc.rw.Unlock()
	}()
	if conn, ok := rc.connMap[t.Id]; ok {
		return conn
	}

	conn, err := mmrpc.Dial(rc.remote.NodeId(), t.NodeId, t.Address, &mmrpc.DialOption{
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
