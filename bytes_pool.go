package xbytes

import (
	"fmt"
	"sync"
)

var DefaultBytesPool = NewBytesPool(64, 1<<30)

type BytesPool interface {
	Get(cap int) []byte
	Put(b []byte)
}

func NewBytesPool(minCap, maxCap int) BytesPool {
	min := roundLog2(minCap)
	max := roundLog2(maxCap)
	if max < min {
		panic(fmt.Sprintf("xbytes.bytesPool: normalize min %v, max %v", min, max))
	}
	sp := &bytesPool{
		min:   min,
		max:   max,
		pools: make([]*sync.Pool, max-min+1),
	}
	return sp
}

type bytesPool struct {
	min   int
	max   int
	pools []*sync.Pool
	sync.RWMutex
}

func (sp *bytesPool) extractPool(cap int) *sync.Pool {
	normalize := roundLog2(cap)
	if normalize > sp.max {
		panic(fmt.Sprintf("xbytes.bytesPool: except <%v, but %v and normalize %v", sp.max, cap, normalize))
	}
	n := normalize - sp.min
	pool := sp.pools[n]
	if pool == nil {
		sp.Lock()
		if pool == nil {
			pool = &sync.Pool{
				New: func() interface{} {
					return make([]byte, 1<<uint(normalize))
				},
			}
			sp.pools[n] = pool
		}
		sp.Unlock()
	}
	return pool
}

func (sp *bytesPool) Get(cap int) []byte {
	pool := sp.extractPool(cap)
	o := pool.Get()
	b := o.([]byte)
	return b[:cap]
}

func (sp *bytesPool) Put(b []byte) {
	if b == nil {
		return
	}
	pool := sp.extractPool(cap(b))
	pool.Put(b)
}
