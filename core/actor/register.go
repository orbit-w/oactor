package actor

import (
	cmap "github.com/orcaman/concurrent-map"
	"github.com/zeebo/xxh3"
	"strconv"
	"sync/atomic"
)

const (
	BucketNum = 64
)

type Register struct {
	id        atomic.Uint64
	engine    *Engine
	localPIDs []cmap.ConcurrentMap
}

func NewRegister(bucketNum int) *Register {
	num := BucketNum
	if bucketNum > 0 {
		num = bucketNum
	}
	re := new(Register)
	re.localPIDs = make([]cmap.ConcurrentMap, num)
	for i := 0; i < num; i++ {
		re.localPIDs[i] = cmap.New()
	}
	return re
}

func (r *Register) Get(pid *PID) (IProcess, bool) {
	b := r.getBucket(pid.Id)
	ref, ok := b.Get(pid.Id)
	return ref.(IProcess), ok
}

func (r *Register) Set(uuid string, p IProcess) (*PID, bool) {
	b := r.getBucket(uuid)
	return &PID{
		Address: r.engine.localAddress(),
		Id:      uuid,
	}, b.SetIfAbsent(uuid, p)
}

func (r *Register) Del(pid *PID) {
	b := r.getBucket(pid.Id)

	ref, _ := b.Pop(pid.Id)
	if actor, ok := ref.(*OActor); ok {
		actor.die()
	}
}

func (r *Register) Id() string {
	v := r.id.Add(1)
	return strconv.FormatUint(v, 10)
}

func (r *Register) getBucket(key string) cmap.ConcurrentMap {
	hash := xxh3.Hash([]byte(key))
	index := int(hash) % len(r.localPIDs)

	return r.localPIDs[index]
}
