package xbytes

import (
	"bytes"
	"fmt"
	"sync"
)

var globalBufferPool = NewBufferPool(16, 1<<26)

func GetBuffer(cap int) *bytes.Buffer {
	return globalBufferPool.GetBuffer(cap)
}

func PutBuffer(buf *bytes.Buffer) {
	globalBufferPool.PutBuffer(buf)
}

type BufferPool interface {
	GetBuffer(cap int) *bytes.Buffer
	PutBuffer(buf *bytes.Buffer)
}

func NewBufferPool(minCap, maxCap int) BufferPool {
	min := roundLog2(minCap)
	max := roundLog2(maxCap)
	if max < min {
		panic(fmt.Sprintf("[Buffer] normalize min %v, max %v", min, max))
	}
	bp := &bufferPool{
		min:   min,
		max:   max,
		pools: make([]*sync.Pool, max-min+1),
	}
	return bp
}

type bufferPool struct {
	min   int
	max   int
	pools []*sync.Pool
	sync.RWMutex
}

func (bp *bufferPool) extractPool(cap int) *sync.Pool {
	normalize := roundLog2(cap)
	if normalize > bp.max {
		panic(fmt.Sprintf("[Buffer] except <%v, but %v and normalize %v", bp.max, cap, normalize))
	}
	n := normalize - bp.min
	pool := bp.pools[n]
	if pool == nil {
		bp.Lock()
		if pool == nil {
			pool = &sync.Pool{
				New: func() interface{} {
					return bytes.NewBuffer(make([]byte, 0, 1<<uint(normalize)))
				},
			}
			bp.pools[n] = pool
		}
		bp.Unlock()
	}
	return pool
}

func (bp *bufferPool) GetBuffer(cap int) *bytes.Buffer {
	pool := bp.extractPool(cap)
	buf := pool.Get()
	return buf.(*bytes.Buffer)
}

func (bp *bufferPool) PutBuffer(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	buf.Reset()
	pool := bp.extractPool(buf.Cap())
	pool.Put(buf)
}
