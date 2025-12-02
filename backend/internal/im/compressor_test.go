package im

import "testing"

func TestCompressor(t *testing.T) {
	c := NewGzipCompressor()
	rawData := []byte("This is a test string to be compressed and decompressed.")

	compressedData, err := c.Compress(rawData)
	if err != nil {
		t.Fatalf("Compression failed: %v", err)
	}
	decompressedData, err := c.Decompress(compressedData)
	if err != nil {
		t.Fatalf("Decompression failed: %v", err)
	}
	if string(decompressedData) != string(rawData) {
		t.Fatalf("Decompressed data does not match original. Got: %s, Want: %s", decompressedData, rawData)
	}
	t.Logf("Original Data: %s", rawData)
	t.Logf("Compressed Data: %v", compressedData)
	t.Logf("Decompressed Data: %s", decompressedData)
}
