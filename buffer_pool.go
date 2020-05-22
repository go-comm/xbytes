package xbytes

import (
	"bytes"
	"fmt"
	"sync"
)

var (
	DefaultBufferPool = NewBufferPool(16, 1<<30)

	DefaultSingleBufferPool = NewSingleBufferPool()
)

type BufferPool interface {
	Get(cap int) *bytes.Buffer
	Put(buf *bytes.Buffer)
}

func NewBufferPool(minCap, maxCap int) BufferPool {
	min := roundLog2(minCap)
	max := roundLog2(maxCap)
	if max < min {
		panic(fmt.Sprintf("xbytes.bufferpool: normalize min %v, max %v", min, max))
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
		panic(fmt.Sprintf("xbytes.bufferpool: except <%v, but %v and normalize %v", bp.max, cap, normalize))
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

func (bp *bufferPool) Get(cap int) *bytes.Buffer {
	pool := bp.extractPool(cap)
	buf := pool.Get()
	return buf.(*bytes.Buffer)
}

func (bp *bufferPool) Put(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	buf.Reset()
	pool := bp.extractPool(buf.Cap())
	pool.Put(buf)
}

type SingleBufferPool interface {
	Get() *bytes.Buffer
	Put(buf *bytes.Buffer)
}

func NewSingleBufferPool() SingleBufferPool {
	return &singleBufferPool{}
}

type singleBufferPool struct {
	pool sync.Pool
}

func (bp *singleBufferPool) Get() *bytes.Buffer {
	v := bp.pool.Get()
	if v != nil {
		return v.(*bytes.Buffer)
	}
	return new(bytes.Buffer)
}

func (bp *singleBufferPool) Put(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	buf.Reset()
	bp.pool.Put(buf)
}
