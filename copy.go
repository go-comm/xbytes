package xbytes

import (
	"io"
	"sync"
)

var zeroCopyPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 4096)
	},
}

func ZeroCopy(dst io.Writer, src io.Reader) (int64, error) {
	buf := zeroCopyPool.Get().([]byte)
	defer zeroCopyPool.Put(buf)
	return io.CopyBuffer(dst, src, buf)
}

var Copy = io.Copy

var CopyN = io.CopyN

var CopyBuffer = io.CopyBuffer
