package xbytes

import (
	"fmt"
	"sync"
)

var (
	DefaultBytebufPool = NewBytebufPool(64, 1<<30)
)

type BytebufPool interface {
	Get(cap int) *Bytebuf
	Put(buf *Bytebuf)
}

func NewBytebufPool(minCap, maxCap int) BytebufPool {
	min := roundLog2(minCap)
	max := roundLog2(maxCap)
	if max < min {
		panic(fmt.Sprintf("xbytes.bytebufpool: normalize min %v, max %v", min, max))
	}
	bp := &bytebufPool{
		min:   min,
		max:   max,
		pools: make([]*sync.Pool, max-min+1),
	}
	return bp
}

type bytebufPool struct {
	min   int
	max   int
	pools []*sync.Pool
	sync.RWMutex
}

func (bp *bytebufPool) extractPool(cap int) *sync.Pool {
	normalize := roundLog2(cap)
	if normalize > bp.max {
		panic(fmt.Sprintf("xbytes.bytebufpool: except <%v, but %v and normalize %v", bp.max, cap, normalize))
	}
	n := normalize - bp.min
	pool := bp.pools[n]
	if pool == nil {
		bp.Lock()
		if pool == nil {
			pool = &sync.Pool{
				New: func() interface{} {
					return NewBytebuf(make([]byte, 0, 1<<uint(normalize)))
				},
			}
			bp.pools[n] = pool
		}
		bp.Unlock()
	}
	return pool
}

func (bp *bytebufPool) Get(cap int) *Bytebuf {
	pool := bp.extractPool(cap)
	buf := pool.Get()
	return buf.(*Bytebuf)
}

func (bp *bytebufPool) Put(buf *Bytebuf) {
	if buf == nil {
		return
	}
	buf.Reset()
	pool := bp.extractPool(buf.Cap())
	pool.Put(buf)
}
