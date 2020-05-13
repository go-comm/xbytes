package xbytes

import "testing"

func Test_Bytebuf_Encode(t *testing.T) {
	var buf Bytebuf

	buf.EncodeVarint(1000)

	buf.EncodeVarint(200)

	t.Log(buf.DecodeVarint())

	t.Log(buf.DecodeVarint())
}
