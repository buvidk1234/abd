package im

import (
	"bytes"
	"compress/gzip"
)

type Compressor interface {
	Compress(rawData []byte) ([]byte, error)
	Decompress(compressedData []byte) ([]byte, error)
}

type GzipCompressor struct{}

func NewGzipCompressor() *GzipCompressor {
	return &GzipCompressor{}
}

func (c *GzipCompressor) Compress(rawData []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(rawData)
	if err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *GzipCompressor) Decompress(compressedData []byte) ([]byte, error) {
	buf := bytes.NewBuffer(compressedData)
	zr, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	var outBuf bytes.Buffer
	if _, err := outBuf.ReadFrom(zr); err != nil {
		return nil, err
	}

	if err := zr.Close(); err != nil {
		return nil, err
	}

	return outBuf.Bytes(), nil
}
