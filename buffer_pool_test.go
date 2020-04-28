package xbytes

import "testing"

func Test_Buffer(t *testing.T) {

	p := NewBufferPool(1022, 1022)

	buf := p.GetBuffer(1024)

	t.Log(buf.Cap())

	p.PutBuffer(buf)

}
