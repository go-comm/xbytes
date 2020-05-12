package xbytes

import (
	"encoding/binary"
	"fmt"
	"io"
)

func NewBytebuf(buf []byte) *Bytebuf {
	return &Bytebuf{buf: buf}
}

type Bytebuf struct {
	off int
	buf []byte
}

func (b *Bytebuf) empty() bool {
	return len(b.buf) <= b.off
}

func (b *Bytebuf) Len() int {
	return len(b.buf) - b.off
}

func (b *Bytebuf) Cap() int {
	return cap(b.buf)
}

func (b *Bytebuf) Limit() int {
	return len(b.buf) - b.off
}

func (b *Bytebuf) Offset() int {
	return b.off
}

func (b *Bytebuf) SetOffset(n int) (err error) {
	if n < 0 || n >= len(b.buf) {
		err = fmt.Errorf("bytebuf: Expected >%d or <%d, not but %d", 0, len(b.buf), n)
	}
	return
}

func (b *Bytebuf) Next(n int) (err error) {
	if err = b.tryRead(n); err != nil {
		return
	}
	b.off += n
	return
}

func (b *Bytebuf) tryRead(n int) error {
	limit := b.Limit()
	if n < 0 || n > limit {
		return fmt.Errorf("bytebuf: Unable to read %d bytes, limit %d", n, limit)
	}
	return nil
}

func (b *Bytebuf) Reset() {
	b.buf = b.buf[:0]
	b.off = 0
}
func (b *Bytebuf) SetBuf(buf []byte) error {
	b.buf = buf
	return nil
}

func (b *Bytebuf) Buf() []byte {
	return b.buf
}

func (b *Bytebuf) Bytes() []byte {
	return b.buf[b.off:]
}

func (b *Bytebuf) Read(p []byte) (int, error) {
	if b.empty() {
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n := copy(p, b.buf[b.off:])
	b.off += n
	return n, nil
}

func (b *Bytebuf) ReadByte() (c byte, err error) {
	if err = b.tryRead(1); err != nil {
		return
	}
	c = b.buf[b.off]
	b.off++
	return
}

func (b *Bytebuf) Write(p []byte) (int, error) {
	b.buf = append(b.buf, p...)
	n := len(p)
	return n, nil
}

func (b *Bytebuf) WriteByte(c byte) error {
	b.buf = append(b.buf, c)
	return nil
}

func (b *Bytebuf) WriteString(s string) (int, error) {
	b.buf = append(b.buf, s[:]...)
	return len(s), nil
}

func (b *Bytebuf) EncodeUint8(x uint8) error {
	b.buf = append(b.buf, x)
	return nil
}

func (b *Bytebuf) EncodeInt8(x int8) error {
	return b.EncodeUint8(uint8(x))
}

func (b *Bytebuf) EncodeUint16(x uint16) error {
	b.buf = append(b.buf,
		uint8(x),
		uint8(x>>8))
	return nil
}

func (b *Bytebuf) EncodeInt16(x int16) error {
	return b.EncodeUint16(uint16(x))
}

func (b *Bytebuf) EncodeUint32(x uint32) error {
	b.buf = append(b.buf,
		uint8(x),
		uint8(x>>8),
		uint8(x>>16),
		uint8(x>>24))
	return nil
}

func (b *Bytebuf) EncodeInt32(x int32) error {
	return b.EncodeUint32(uint32(x))
}

func (b *Bytebuf) EncodeUint64(x uint64) error {
	b.buf = append(b.buf,
		uint8(x),
		uint8(x>>8),
		uint8(x>>16),
		uint8(x>>24),
		uint8(x>>32),
		uint8(x>>40),
		uint8(x>>48),
		uint8(x>>56))
	return nil
}

func (b *Bytebuf) EncodeInt64(x uint64) error {
	return b.EncodeUint64(uint64(x))
}

func (b *Bytebuf) EncodeUvarint(x uint64) error {
	for x >= 1<<7 {
		b.buf = append(b.buf, uint8(x&0x7f|0x80))
		x >>= 7
	}
	b.buf = append(b.buf, uint8(x))
	return nil
}

func (b *Bytebuf) EncodeVarint(x int64) error {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return b.EncodeUvarint(ux)
}

func (b *Bytebuf) EncodeBytes(p []byte) error {
	b.EncodeUvarint(uint64(len(p)))
	b.Write(p)
	return nil
}

func (b *Bytebuf) EncodeString(s string) error {
	b.EncodeUvarint(uint64(len(s)))
	b.Write(StringToBytes(&s))
	return nil
}

func (b *Bytebuf) DecodeUint8() (x uint8, err error) {
	if err = b.tryRead(1); err != nil {
		return
	}
	i := b.off + 1
	b.off = i
	x = uint8(b.buf[i])
	return
}

func (b *Bytebuf) DecodeInt8() (x int8, err error) {
	ux, err := b.DecodeUint8()
	return int8(ux), err
}

func (b *Bytebuf) DecodeUint16() (x uint16, err error) {
	if err = b.tryRead(1); err != nil {
		return
	}
	i := b.off + 2
	b.off = i
	x = uint16(b.buf[i-2])
	x |= uint16(b.buf[i-1]) << 8
	return
}

func (b *Bytebuf) DecodeInt16() (x int16, err error) {
	ux, err := b.DecodeUint16()
	return int16(ux), err
}

func (b *Bytebuf) DecodeUint32() (x uint32, err error) {
	if err = b.tryRead(1); err != nil {
		return
	}
	i := b.off + 4
	b.off = i
	x = uint32(b.buf[i-4])
	x |= uint32(b.buf[i-3]) << 8
	x |= uint32(b.buf[i-2]) << 16
	x |= uint32(b.buf[i-1]) << 24
	return
}

func (b *Bytebuf) DecodeInt32() (x int32, err error) {
	ux, err := b.DecodeUint32()
	return int32(ux), err
}

func (b *Bytebuf) DecodeUint64() (x uint64, err error) {
	if err = b.tryRead(1); err != nil {
		return
	}
	i := b.off + 8
	b.off = i
	x = uint64(b.buf[i-8])
	x |= uint64(b.buf[i-7]) << 8
	x |= uint64(b.buf[i-6]) << 16
	x |= uint64(b.buf[i-5]) << 24
	x |= uint64(b.buf[i-4]) << 32
	x |= uint64(b.buf[i-3]) << 40
	x |= uint64(b.buf[i-2]) << 48
	x |= uint64(b.buf[i-1]) << 56
	return
}

func (b *Bytebuf) DecodeInt64() (x int64, err error) {
	ux, err := b.DecodeUint64()
	return int64(ux), err
}

func (b *Bytebuf) DecodeUvarint() (x uint64, err error) {
	return binary.ReadUvarint(b)
}

func (b *Bytebuf) DecodeVarint() (x int64, err error) {
	return binary.ReadVarint(b)
}

func (b *Bytebuf) DecodeBytes(alloc bool) (p []byte, err error) {
	var n uint64
	n, err = b.DecodeUvarint()
	if err != nil {
		return nil, err
	}
	nb := int(n)
	if err = b.tryRead(nb); err != nil {
		return nil, err
	}
	if !alloc {
		p = b.buf[b.off : b.off+nb]
		b.off += nb
		return
	}
	p = make([]byte, nb)
	copy(p, b.buf[b.off:])
	b.off += nb
	return
}

func (b *Bytebuf) DecodeString() (string, error) {
	p, err := b.DecodeBytes(false)
	if err != nil {
		return "", err
	}
	return string(p), nil
}
