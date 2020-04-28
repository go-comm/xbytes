package xbytes

import (
	"bytes"
	"io"
)

func NewBuffer(b []byte) *bytes.Buffer {
	return bytes.NewBuffer(b)
}

func NewReaderFromString(s *string) io.Reader {
	return bytes.NewReader(StringToBytes(s))
}

func NewBufferFromString(s *string) *bytes.Buffer {
	return bytes.NewBuffer(StringToBytes(s))
}

func BufferToString(b *bytes.Buffer) string {
	return BytesToString(b.Bytes())
}
