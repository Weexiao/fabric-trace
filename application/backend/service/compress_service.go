package service

import (
	"backend/pkg"
	"encoding/json"
	"fmt"
	"log"
)

// Compressor 定义数据压缩器的策略接口。
// 短期实现使用 Gzip；长期可通过 BTAECompressor 替换为神经网络压缩。
// Compress 返回 Base64 编码的压缩数据以及链上存证用的 CompressionEvidence。
// Decompress 接受 Base64 编码的压缩字符串，返回原始字节。
type Compressor interface {
	// Algorithm 返回压缩算法标识符
	Algorithm() string
	// Compress 压缩原始数据，返回 Base64 编码结果及防篡改证据
	Compress(data []byte) (b64 string, evidence pkg.CompressionEvidence, err error)
	// Decompress 从 Base64 字符串解压还原原始数据
	Decompress(b64 string) ([]byte, error)
}

// ---------- GzipCompressor ----------

// GzipCompressor 基于 Gzip (RFC 1952) 的默认压缩策略。
// 适用于通用文本/JSON 数据，压缩率典型在 60-85%。
type GzipCompressor struct{}

func (g *GzipCompressor) Algorithm() string { return "gzip" }

func (g *GzipCompressor) Compress(data []byte) (string, pkg.CompressionEvidence, error) {
	compressed, err := pkg.GzipCompress(data)
	if err != nil {
		return "", pkg.CompressionEvidence{}, fmt.Errorf("gzip compress: %w", err)
	}
	evidence := pkg.BuildCompressionEvidence("gzip", data, compressed)
	b64 := pkg.Base64Encode(compressed)
	return b64, evidence, nil
}

func (g *GzipCompressor) Decompress(b64 string) ([]byte, error) {
	compressed, err := pkg.Base64Decode(b64)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}
	return pkg.GzipDecompress(compressed)
}

// ---------- BTAECompressor (预留) ----------

// BTAECompressor 为 Bernoulli Transformer Autoencoder 预留的压缩策略。
// 当前实现降级为 GzipCompressor，当模型服务就绪后可替换为 gRPC/HTTP 调用推理服务。
//
// 技术原理：BTAE 将输入数据映射到低维 Bernoulli 潜空间，生成定长二进制特征向量。
// 该特征向量同时作为压缩表示和异常检测的依据——重构误差超过阈值即标记异常。
// FeatureVector 字段将存储该潜空间表示，支持后续链上异常检测扩展。
type BTAECompressor struct {
	ModelPath string // 预留：ONNX/TorchScript 模型路径或推理服务地址
	fallback  *GzipCompressor
}

func NewBTAECompressor(modelPath string) *BTAECompressor {
	return &BTAECompressor{
		ModelPath: modelPath,
		fallback:  &GzipCompressor{},
	}
}

func (b *BTAECompressor) Algorithm() string { return "btae" }

func (b *BTAECompressor) Compress(data []byte) (string, pkg.CompressionEvidence, error) {
	// TODO: 当 BTAE 推理服务就绪后，替换以下逻辑：
	//   1. 将 data 发送至 BTAE 推理服务 (gRPC/HTTP)
	//   2. 获取定长特征向量 (Latent Representation)
	//   3. 将特征向量填入 evidence.FeatureVector
	//   4. 使用 Bernoulli 编码后的二进制流作为压缩结果
	log.Printf("[BTAECompressor] BTAE model not yet deployed (path=%s), falling back to gzip", b.ModelPath)

	b64, evidence, err := b.fallback.Compress(data)
	if err != nil {
		return "", pkg.CompressionEvidence{}, err
	}
	// 标记算法为 btae_fallback 以便区分
	evidence.Algorithm = "btae_fallback_gzip"
	// 预留：evidence.FeatureVector = btaeInfer(data)
	return b64, evidence, nil
}

func (b *BTAECompressor) Decompress(b64 string) ([]byte, error) {
	// BTAE 解码需要模型推理；当前降级使用 Gzip
	log.Printf("[BTAECompressor] BTAE decode not available, using gzip fallback")
	return b.fallback.Decompress(b64)
}

// ---------- Factory ----------

// NewCompressor 根据算法名称创建对应的 Compressor 实例。
// 支持的算法: "gzip" (默认), "btae"。
// 在 config.yaml 中配置 compression.algorithm 即可切换策略。
func NewCompressor(algorithm string) Compressor {
	switch algorithm {
	case "btae":
		return NewBTAECompressor("")
	default:
		return &GzipCompressor{}
	}
}

// ---------- 辅助函数 ----------

// CompressJSON 将任意结构体序列化为 JSON 后进行压缩，返回 Base64 与证据。
// 适用于上链前的批量数据压缩场景。
func CompressJSON(c Compressor, v interface{}) (string, pkg.CompressionEvidence, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", pkg.CompressionEvidence{}, fmt.Errorf("json marshal: %w", err)
	}
	return c.Compress(data)
}

// DecompressJSON 从 Base64 解压后反序列化为目标结构体。
func DecompressJSON(c Compressor, b64 string, target interface{}) error {
	data, err := c.Decompress(b64)
	if err != nil {
		return fmt.Errorf("decompress: %w", err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}
	return nil
}
