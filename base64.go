package xbytes

import (
	"bytes"
	"encoding/base64"
	"io"
)

var (
	Base64Std = NewBase64Stream(base64.StdEncoding)
	Base64URL = NewBase64Stream(base64.URLEncoding)
)

func NewBase64Stream(en *base64.Encoding) *base64Stream {
	return &base64Stream{en: en}
}

type base64Stream struct {
	en *base64.Encoding
}

func (e *base64Stream) EncodeFromReader(w io.Writer, r io.Reader) error {
	ew := base64.NewEncoder(e.en, w)
	ZeroCopy(ew, r)
	return ew.Close()
}

func (e *base64Stream) Encode(w io.Writer, data []byte) error {
	ew := base64.NewEncoder(e.en, w)
	ew.Write(data)
	return ew.Close()
}

func (e *base64Stream) DecodeFromReader(w io.Writer, r io.Reader) error {
	er := base64.NewDecoder(e.en, r)
	_, err := ZeroCopy(w, er)
	return err
}

func (e *base64Stream) Decode(w io.Writer, data []byte) error {
	return e.DecodeFromReader(w, bytes.NewReader(data))
}

func (e *base64Stream) DecodeToBytes(s *string) ([]byte, error) {
	dbuf := make([]byte, e.en.DecodedLen(len(*s)))
	n, err := e.en.Decode(dbuf, StringToBytes(s))
	return dbuf[:n], err
}
