package pkg

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// CompressionEvidence 记录压缩操作的防篡改证据，用于链上存证。
// 包含原始数据与压缩数据的 SHA-256 摘要、字节大小及压缩率，
// 使得链上节点可独立验证数据完整性，同时为后续 BTAE 异常检测预留 Algorithm 字段。
type CompressionEvidence struct {
	Algorithm        string  `json:"algorithm"`        // 压缩算法标识: "gzip", "zlib", "btae"
	OriginalHash     string  `json:"originalHash"`     // 原始数据 SHA-256 hex
	CompressedHash   string  `json:"compressedHash"`   // 压缩后数据 SHA-256 hex
	OriginalSize     int64   `json:"originalSize"`     // 原始数据字节数
	CompressedSize   int64   `json:"compressedSize"`   // 压缩后字节数
	CompressionRatio float64 `json:"compressionRatio"` // 压缩率 = compressedSize / originalSize
	// FeatureVector 为 BTAE 模型预留：定长特征摘要（Latent Representation）
	FeatureVector []float64 `json:"featureVector,omitempty"`
}

// GzipCompress 使用 Gzip (RFC 1952) 对数据进行压缩，返回压缩后的字节流。
// 采用 BestCompression 级别以最大化压缩比。
func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("gzip writer init failed: %w", err)
	}
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("gzip write failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("gzip close failed: %w", err)
	}
	return buf.Bytes(), nil
}

// GzipDecompress 对 Gzip 压缩的数据进行解压缩。
func GzipDecompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("gzip reader init failed: %w", err)
	}
	defer r.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("gzip read failed: %w", err)
	}
	return out, nil
}

// Base64Encode 将字节流编码为 URL-safe Base64 字符串（无填充），
// 适合在 JSON 传输中嵌入二进制数据。
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode 将 Base64 字符串解码为字节流。
func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// sha256Hex 计算数据的 SHA-256 摘要，返回 hex 字符串。
func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// BuildCompressionEvidence 根据原始数据和压缩后的数据构建完整的压缩证据。
// 该证据将被序列化后存储到区块链上，用于防篡改验证。
func BuildCompressionEvidence(algorithm string, original, compressed []byte) CompressionEvidence {
	origSize := int64(len(original))
	compSize := int64(len(compressed))
	var ratio float64
	if origSize > 0 {
		ratio = float64(compSize) / float64(origSize)
	}
	return CompressionEvidence{
		Algorithm:        algorithm,
		OriginalHash:     sha256Hex(original),
		CompressedHash:   sha256Hex(compressed),
		OriginalSize:     origSize,
		CompressedSize:   compSize,
		CompressionRatio: ratio,
	}
}
