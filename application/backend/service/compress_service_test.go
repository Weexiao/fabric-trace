package service

import (
	"encoding/json"
	"testing"
)

func TestGzipCompressorRoundTrip(t *testing.T) {
	c := NewCompressor("gzip")
	original := []byte(`{"arg1":"苹果","arg2":"四川省","arg3":"2026-01-01","arg4":"2026-01-05","arg5":"张三"}`)

	b64, evidence, err := c.Compress(original)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}
	if b64 == "" {
		t.Fatal("base64 result is empty")
	}
	if evidence.Algorithm != "gzip" {
		t.Errorf("algorithm=%s, want gzip", evidence.Algorithm)
	}
	t.Logf("compressed base64 length=%d, evidence=%+v", len(b64), evidence)

	restored, err := c.Decompress(b64)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}
	if string(restored) != string(original) {
		t.Fatalf("round-trip mismatch")
	}
}

func TestBTAECompressorFallback(t *testing.T) {
	c := NewCompressor("btae")
	data := []byte(`{"test":"btae fallback"}`)

	b64, evidence, err := c.Compress(data)
	if err != nil {
		t.Fatalf("BTAE Compress failed: %v", err)
	}
	if evidence.Algorithm != "btae_fallback_gzip" {
		t.Errorf("algorithm=%s, want btae_fallback_gzip", evidence.Algorithm)
	}

	restored, err := c.Decompress(b64)
	if err != nil {
		t.Fatalf("BTAE Decompress failed: %v", err)
	}
	if string(restored) != string(data) {
		t.Fatalf("round-trip mismatch")
	}
}

func TestCompressJSON(t *testing.T) {
	type payload struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	c := NewCompressor("gzip")
	p := payload{Name: "测试", Value: 42}

	b64, evidence, err := CompressJSON(c, p)
	if err != nil {
		t.Fatalf("CompressJSON failed: %v", err)
	}
	if evidence.OriginalSize == 0 {
		t.Error("original size should not be 0")
	}

	var restored payload
	if err := DecompressJSON(c, b64, &restored); err != nil {
		t.Fatalf("DecompressJSON failed: %v", err)
	}
	if restored.Name != p.Name || restored.Value != p.Value {
		t.Fatalf("mismatch: got %+v, want %+v", restored, p)
	}

	// Verify evidence can be serialized to JSON (for chaincode)
	eb, err := json.Marshal(evidence)
	if err != nil {
		t.Fatalf("json.Marshal evidence failed: %v", err)
	}
	t.Logf("evidence JSON: %s", string(eb))
}

func TestNewCompressorDefault(t *testing.T) {
	c := NewCompressor("unknown")
	if c.Algorithm() != "gzip" {
		t.Errorf("default compressor should be gzip, got %s", c.Algorithm())
	}
}
