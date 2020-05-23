package xbytes

import (
	"encoding/binary"
	"io"
)

func NewBytesWriter(w io.Writer) *BytesWriter {
	if bw, ok := w.(*BytesWriter); ok {
		return bw
	}
	if bwriter, ok := w.(io.ByteWriter); ok {
		return &BytesWriter{nil, w, bwriter}
	}
	return &BytesWriter{nil, w, nil}
}

type BytesWriter struct {
	err error
	io.Writer
	bwriter io.ByteWriter
}

func (bw *BytesWriter) Write(p []byte) (n int, err error) {
	if bw.err != nil {
		err = bw.err
		return
	}
	n, err = bw.Writer.Write(p)
	bw.err = err
	return
}

func (bw *BytesWriter) WriteByte(c byte) (err error) {
	if bw.bwriter != nil {
		err = bw.bwriter.WriteByte(c)
		bw.err = err
		return
	}
	p := []byte{c}
	_, err = bw.Writer.Write(p)
	bw.err = err
	return
}

func (bw *BytesWriter) EncodeUint8(v uint8) (err error) {
	return bw.WriteByte(byte(v))
}

func (bw *BytesWriter) EncodeInt8(v int8) (err error) {
	return bw.WriteByte(byte(v))
}

func (bw *BytesWriter) EncodeUint16(v uint16) (err error) {
	p := []byte{0, 0}
	binary.BigEndian.PutUint16(p, v)
	_, err = bw.Write(p)
	return
}

func (bw *BytesWriter) EncodeInt16(v int16) (err error) {
	return bw.EncodeUint16(uint16(v))
}

func (bw *BytesWriter) EncodeUint32(v uint32) (err error) {
	p := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(p, v)
	_, err = bw.Write(p)
	return
}

func (bw *BytesWriter) EncodeInt32(v int32) (err error) {
	return bw.EncodeUint32(uint32(v))
}

func (bw *BytesWriter) EncodeInt(v int) (err error) {
	return bw.EncodeUint32(uint32(v))
}

func (bw *BytesWriter) EncodeUint64(v uint64) (err error) {
	p := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint64(p, v)
	_, err = bw.Write(p)
	return
}

func (bw *BytesWriter) EncodeInt64(v int64) (err error) {
	return bw.EncodeUint64(uint64(v))
}

func (bw *BytesWriter) EncodeUvarint(v uint64) (n int, err error) {
	p := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	n = binary.PutUvarint(p, v)
	_, err = bw.Write(p[:n])
	return
}

func (bw *BytesWriter) EncodeVarint(v int64) (n int, err error) {
	p := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	n = binary.PutVarint(p, v)
	_, err = bw.Write(p[:n])
	return
}

func (bw *BytesWriter) EncodeBytes(p []byte) (n int, err error) {
	var nm int
	n, err = bw.EncodeUvarint(uint64(len(p)))
	if err != nil {
		return 0, err
	}
	nm, err = bw.Write(p)
	n += nm
	return
}

func (bw *BytesWriter) EncodeString(s string) (n int, err error) {
	return bw.EncodeBytes(StringToBytes(&s))
}
