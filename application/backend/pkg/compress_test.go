package pkg

import (
	"testing"
)

func TestGzipRoundTrip(t *testing.T) {
	original := []byte(`{"arg1":"苹果","arg2":"四川省成都市","arg3":"2026-01-01 08:00:00","arg4":"2026-01-05 10:00:00","arg5":"张三果园"}`)

	compressed, err := GzipCompress(original)
	if err != nil {
		t.Fatalf("GzipCompress failed: %v", err)
	}
	if len(compressed) == 0 {
		t.Fatal("compressed data is empty")
	}
	t.Logf("original=%d bytes, compressed=%d bytes, ratio=%.2f%%",
		len(original), len(compressed), float64(len(compressed))/float64(len(original))*100)

	decompressed, err := GzipDecompress(compressed)
	if err != nil {
		t.Fatalf("GzipDecompress failed: %v", err)
	}
	if string(decompressed) != string(original) {
		t.Fatalf("round-trip mismatch:\n  original:     %s\n  decompressed: %s", original, decompressed)
	}
}

func TestBase64RoundTrip(t *testing.T) {
	data := []byte("hello, 区块链压缩测试")
	encoded := Base64Encode(data)
	decoded, err := Base64Decode(encoded)
	if err != nil {
		t.Fatalf("Base64Decode failed: %v", err)
	}
	if string(decoded) != string(data) {
		t.Fatalf("round-trip mismatch")
	}
}

func TestBuildCompressionEvidence(t *testing.T) {
	original := []byte("test data for compression evidence")
	compressed, _ := GzipCompress(original)

	ev := BuildCompressionEvidence("gzip", original, compressed)

	if ev.Algorithm != "gzip" {
		t.Errorf("algorithm=%s, want gzip", ev.Algorithm)
	}
	if ev.OriginalSize != int64(len(original)) {
		t.Errorf("originalSize=%d, want %d", ev.OriginalSize, len(original))
	}
	if ev.CompressedSize != int64(len(compressed)) {
		t.Errorf("compressedSize=%d, want %d", ev.CompressedSize, len(compressed))
	}
	if ev.OriginalHash == "" || ev.CompressedHash == "" {
		t.Error("hash should not be empty")
	}
	if ev.CompressionRatio <= 0 {
		t.Errorf("unexpected compression ratio: %.4f", ev.CompressionRatio)
	}
	t.Logf("evidence: algo=%s, origHash=%s..., compHash=%s..., ratio=%.4f",
		ev.Algorithm, ev.OriginalHash[:16], ev.CompressedHash[:16], ev.CompressionRatio)
}

func TestGzipCompressEmptyData(t *testing.T) {
	compressed, err := GzipCompress([]byte{})
	if err != nil {
		t.Fatalf("GzipCompress on empty data failed: %v", err)
	}
	decompressed, err := GzipDecompress(compressed)
	if err != nil {
		t.Fatalf("GzipDecompress on empty data failed: %v", err)
	}
	if len(decompressed) != 0 {
		t.Fatalf("expected empty, got %d bytes", len(decompressed))
	}
}
