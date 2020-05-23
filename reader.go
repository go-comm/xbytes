package xbytes

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	maxBytesLength = (1 << 16) - 1
)

var (
	ErrLargeBytesLength = errors.New("xbytes.BytesReader: large bytes length")
)

func NewBytesReader(r io.Reader) *BytesReader {
	if br, ok := r.(*BytesReader); ok {
		return br
	}
	if breader, ok := r.(io.ByteReader); ok {
		return &BytesReader{nil, r, breader}
	}
	return &BytesReader{nil, r, nil}
}

type BytesReader struct {
	err error
	io.Reader
	breader io.ByteReader
}

func (br *BytesReader) Error() error {
	return br.err
}

func (br *BytesReader) ReadAtLeast(r io.Reader, buf []byte, min int) (n int, err error) {
	if br.err != nil {
		return 0, br.err
	}
	n, err = io.ReadAtLeast(r, buf, min)
	br.err = err
	return
}

func (br *BytesReader) Read(p []byte) (n int, err error) {
	if br.err != nil {
		return 0, br.err
	}
	n, err = br.Reader.Read(p)
	br.err = err
	return
}

func (br *BytesReader) ReadByte() (byte, error) {
	if br.breader != nil {
		v, err := br.breader.ReadByte()
		br.err = err
		return v, err
	}
	b := []byte{0}
	n, err := br.Read(b)
	if err != nil {
		return 0, err
	}
	if n < len(b) {
		return 0, io.ErrUnexpectedEOF
	}
	return b[0], nil
}

func (br *BytesReader) DecodeUint8() (uint8, error) {
	v, err := br.ReadByte()
	return uint8(v), err
}

func (br *BytesReader) DecodeInt8() (int8, error) {
	v, err := br.ReadByte()
	return int8(v), err
}

func (br *BytesReader) DecodeUint16() (uint16, error) {
	b := []byte{0, 0}
	_, err := br.ReadAtLeast(br, b, 2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(b), nil
}

func (br *BytesReader) DecodeInt16() (int16, error) {
	v, err := br.DecodeUint16()
	return int16(v), err
}

func (br *BytesReader) DecodeUint32() (uint32, error) {
	b := []byte{0, 0, 0, 0}
	_, err := br.ReadAtLeast(br, b, 4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}

func (br *BytesReader) DecodeInt32() (int32, error) {
	v, err := br.DecodeUint16()
	return int32(v), err
}

func (br *BytesReader) DecodeUint64() (uint64, error) {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	_, err := br.ReadAtLeast(br, b, 8)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b), nil
}

func (br *BytesReader) DecodeInt64() (int64, error) {
	v, err := br.DecodeUint64()
	return int64(v), err
}

func (br *BytesReader) DecodeVarint() (int64, error) {
	return binary.ReadVarint(br)
}

func (br *BytesReader) DecodeUvarint() (uint64, error) {
	return binary.ReadUvarint(br)
}

func (br *BytesReader) DecodeBytes(p []byte) (n int, err error) {
	var un uint64
	un, err = br.DecodeUvarint()
	if err != nil {
		return 0, err
	}
	n = int(un)
	return br.ReadAtLeast(br, p, n)
}

func (br *BytesReader) DecodeAllocBytes() (p []byte, err error) {
	var un uint64
	un, err = br.DecodeUvarint()
	if err != nil {
		return nil, err
	}
	n := int(un)
	if n > maxBytesLength {
		br.err = ErrLargeBytesLength
		return nil, ErrLargeBytesLength
	}
	p = make([]byte, n)
	_, err = br.ReadAtLeast(br, p[:], n)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (br *BytesReader) DecodeString() (s string, err error) {
	b, err := br.DecodeAllocBytes()
	if err != nil {
		return "", err
	}
	return BytesToString(b), nil
}
