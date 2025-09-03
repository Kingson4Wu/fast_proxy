package gzip

import (
	"bytes"
	"compress/gzip"
	"github.com/Kingson4Wu/fast_proxy/common/compress/interfaces"
	"io"
)

type Gzip struct {
}

var _ interfaces.Compress = &Gzip{}

func (*Gzip) Encode(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (*Gzip) Decode(compressed []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(compressed))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
